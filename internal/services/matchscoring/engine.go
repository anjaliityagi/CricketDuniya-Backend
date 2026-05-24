package matchscoring

import (
	"CricketDuniya-Backend/internal/database"
	"CricketDuniya-Backend/internal/dto"
	"strings"

	"github.com/jmoiron/sqlx"
)

type Update struct {
	InningsRuns    int `json:"innings_runs"`
	InningsWickets int `json:"innings_wickets"`
}

type Engine struct{}

func NewEngine() *Engine {
	return &Engine{}
}

func (e *Engine) Process(req dto.BallRequest) (*Update, error) {
	tx, err := database.DB.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	wicketInc := 0
	if req.IsWicket {
		wicketInc = 1
	}

	var inningsRuns, inningsWickets int
	err = tx.QueryRowx(`UPDATE innings
		SET
			total_runs = COALESCE(total_runs, 0) + $2,
			total_wickets = COALESCE(total_wickets, 0) + $3
		WHERE id = $1
		RETURNING total_runs, total_wickets`,
		req.InningsID, req.TotalRuns, wicketInc).Scan(&inningsRuns, &inningsWickets)
	if err != nil {
		return nil, err
	}

	legalBall := !isExtraBall(req.BallType)

	if req.StrikerID.String() != "" {
		if err := upsertBattingStats(tx, req.MatchID.String(), req.StrikerID.String(), req.RunsOffBat, legalBall, req.IsBoundaryFour, req.IsBoundarySix, req.IsWicket); err != nil {
			return nil, err
		}
	}

	if req.BowlerID.String() != "" {
		conceded := req.TotalRuns - req.Byes - req.LegByes
		if conceded < 0 {
			conceded = 0
		}
		if err := upsertBowlingStats(tx, req.MatchID.String(), req.BowlerID.String(), conceded, legalBall, req.IsWicket); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &Update{
		InningsRuns:    inningsRuns,
		InningsWickets: inningsWickets,
	}, nil
}

func isExtraBall(ballType string) bool {
	bt := strings.ToLower(strings.TrimSpace(ballType))
	return bt == "wide" || bt == "no_ball"
}

func upsertBattingStats(tx *sqlx.Tx, matchID, matchPlayerID string, runs int, legalBall, isFour, isSix, isOut bool) error {
	balls := 0
	if legalBall {
		balls = 1
	}
	fours := 0
	if isFour {
		fours = 1
	}
	sixes := 0
	if isSix {
		sixes = 1
	}
	outInc := false
	if isOut {
		outInc = true
	}

	res, err := tx.Exec(`
		UPDATE player_match_stats
		SET
			runs_scored = COALESCE(runs_scored, 0) + $3,
			balls_faced = COALESCE(balls_faced, 0) + $4,
			fours = COALESCE(fours, 0) + $5,
			sixes = COALESCE(sixes, 0) + $6,
			is_out = COALESCE(is_out, FALSE) OR $7,
			updated_at = NOW()
		WHERE match_id = $1
		  AND team_player_id = $2
	`, matchID, matchPlayerID, runs, balls, fours, sixes, outInc)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows > 0 {
		return nil
	}

	_, err = tx.Exec(`
		INSERT INTO player_match_stats (
			match_id,
			player_id,
			team_player_id,
			runs_scored,
			balls_faced,
			fours,
			sixes,
			is_out,
			updated_at
		)
		SELECT
			$1,
			tp.player_id,
			tp.id,
			$3,
			$4,
			$5,
			$6,
			$7,
			NOW()
		FROM team_players tp
		WHERE tp.id = $2
		  AND tp.player_id IS NOT NULL
	`, matchID, matchPlayerID, runs, balls, fours, sixes, outInc)
	return err
}

func upsertBowlingStats(tx *sqlx.Tx, matchID, matchPlayerID string, runsConceded int, legalBall, isWicket bool) error {
	wickets := 0
	if isWicket {
		wickets = 1
	}
	legalBallDelta := 0
	if legalBall {
		legalBallDelta = 1
	}

	res, err := tx.Exec(`
		UPDATE player_match_stats
		SET
			runs_conceded = COALESCE(runs_conceded, 0) + $3,
			wickets_taken = COALESCE(wickets_taken, 0) + $4,
			legal_balls_bowled = COALESCE(legal_balls_bowled, 0) + $5::INTEGER,
			overs_bowled = FLOOR((COALESCE(legal_balls_bowled, 0) + $5::INTEGER) / 6.0)
				+ MOD(COALESCE(legal_balls_bowled, 0) + $5::INTEGER, 6)::NUMERIC / 10.0,
			updated_at = NOW()
		WHERE match_id = $1
		  AND team_player_id = $2
	`, matchID, matchPlayerID, runsConceded, wickets, legalBallDelta)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows > 0 {
		return nil
	}

	_, err = tx.Exec(`
		INSERT INTO player_match_stats (
			match_id,
			player_id,
			team_player_id,
			runs_conceded,
			wickets_taken,
			legal_balls_bowled,
			overs_bowled,
			updated_at
		)
		SELECT
			$1,
			tp.player_id,
			tp.id,
			$3,
			$4,
			$5::INTEGER,
			FLOOR($5::INTEGER / 6.0) + MOD($5::INTEGER, 6)::NUMERIC / 10.0,
			NOW()
		FROM team_players tp
		WHERE tp.id = $2
		  AND tp.player_id IS NOT NULL
	`, matchID, matchPlayerID, runsConceded, wickets, legalBallDelta)
	return err
}
