package repositories

import (
	"CricketDuniya-Backend/internal/database"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type BatterTotals struct {
	RunsBefore  int `db:"runs_before"`
	BallsBefore int `db:"balls_before"`
}

type BowlerTotals struct {
	LegalBallsBefore   int `db:"legal_balls_before"`
	RunsConcededBefore int `db:"runs_conceded_before"`
}

type BowlerOverTotals struct {
	OverLegalBefore      int `db:"over_legal_before"`
	OverRunsBefore       int `db:"over_runs_before"`
	OverBoundariesBefore int `db:"over_boundaries_before"`
}

func BeginTx() (*sqlx.Tx, error) { return database.DB.Beginx() }

func GetBatterTotals(matchID, inningsID, strikerID string) (*BatterTotals, error) {
	var out BatterTotals
	err := database.DB.Get(&out, `
		SELECT
			COALESCE(SUM(CASE WHEN ball_type IN ('wide','bye','leg_bye') THEN 0 ELSE GREATEST(total_runs - extras, 0) END), 0) AS runs_before,
			COALESCE(SUM(CASE WHEN ball_type NOT IN ('wide','no_ball','dead_ball','retired_hurt') THEN 1 ELSE 0 END), 0) AS balls_before
		FROM ball_events
		WHERE match_id = $1 AND innings_id = $2 AND striker_id = $3 AND is_deleted = FALSE
	`, matchID, inningsID, strikerID)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func GetBowlerTotals(matchID, inningsID, bowlerID string) (*BowlerTotals, error) {
	var out BowlerTotals
	err := database.DB.Get(&out, `
		SELECT
			COALESCE(SUM(CASE WHEN ball_type NOT IN ('wide','no_ball','dead_ball','retired_hurt') THEN 1 ELSE 0 END), 0) AS legal_balls_before,
			COALESCE(SUM(total_runs - byes - leg_byes), 0) AS runs_conceded_before
		FROM ball_events
		WHERE match_id = $1 AND innings_id = $2 AND bowler_id = $3 AND is_deleted = FALSE
	`, matchID, inningsID, bowlerID)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func GetBowlerOverTotals(matchID, inningsID, bowlerID string, ballNo int) (*BowlerOverTotals, error) {
	var out BowlerOverTotals
	err := database.DB.Get(&out, `
		SELECT
			COALESCE(SUM(CASE WHEN ball_type NOT IN ('wide','no_ball','dead_ball','retired_hurt') THEN 1 ELSE 0 END), 0) AS over_legal_before,
			COALESCE(SUM(total_runs - byes - leg_byes), 0) AS over_runs_before,
			COALESCE(SUM(CASE WHEN ball_type IN ('wide','bye','leg_bye') THEN 0 WHEN GREATEST(total_runs - extras, 0) IN (4, 6) THEN 1 ELSE 0 END), 0) AS over_boundaries_before
		FROM ball_events
		WHERE match_id = $1 AND innings_id = $2 AND bowler_id = $3 AND ball_no = $4 AND is_deleted = FALSE
	`, matchID, inningsID, bowlerID, ballNo)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func ListRecentRunsOffBatForStriker(matchID, inningsID, strikerID string, limit int) ([]int, error) {
	rows, err := database.DB.Queryx(`
		SELECT CASE WHEN ball_type IN ('wide','bye','leg_bye') THEN 0 ELSE GREATEST(total_runs - extras, 0) END AS runs_off_bat
		FROM ball_events
		WHERE match_id = $1 AND innings_id = $2 AND striker_id = $3 AND is_deleted = FALSE
		ORDER BY created_at DESC
		LIMIT $4
	`, matchID, inningsID, strikerID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]int, 0, limit)
	for rows.Next() {
		var v int
		if err := rows.Scan(&v); err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, rows.Err()
}

func ListRecentRunsOffBatForBowler(matchID, inningsID, bowlerID string, limit int) ([]int, error) {
	rows, err := database.DB.Queryx(`
		SELECT CASE WHEN ball_type IN ('wide','bye','leg_bye') THEN 0 ELSE GREATEST(total_runs - extras, 0) END AS runs_off_bat
		FROM ball_events
		WHERE match_id = $1 AND innings_id = $2 AND bowler_id = $3 AND is_deleted = FALSE
		ORDER BY created_at DESC
		LIMIT $4
	`, matchID, inningsID, bowlerID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]int, 0, limit)
	for rows.Next() {
		var v int
		if err := rows.Scan(&v); err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, rows.Err()
}

func ListRecentBallTypesForBowler(matchID, inningsID, bowlerID string, limit int) ([]string, error) {
	rows, err := database.DB.Queryx(`
		SELECT ball_type
		FROM ball_events
		WHERE match_id = $1 AND innings_id = $2 AND bowler_id = $3 AND is_deleted = FALSE
		ORDER BY created_at DESC
		LIMIT $4
	`, matchID, inningsID, bowlerID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]string, 0, limit)
	for rows.Next() {
		var v string
		if err := rows.Scan(&v); err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, rows.Err()
}

func ListRecentWicketsForBowlerInOver(matchID, inningsID, bowlerID string, ballNo, limit int) ([]bool, error) {
	rows, err := database.DB.Queryx(`
		SELECT is_wicket
		FROM ball_events
		WHERE match_id = $1 AND innings_id = $2 AND bowler_id = $3 AND ball_no = $4 AND is_deleted = FALSE
		ORDER BY created_at DESC
		LIMIT $5
	`, matchID, inningsID, bowlerID, ballNo, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]bool, 0, limit)
	for rows.Next() {
		var v bool
		if err := rows.Scan(&v); err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, rows.Err()
}

func UpdateInningsTotalsTx(tx *sqlx.Tx, inningsID string, runsDelta, wicketsDelta int) (int, int, error) {
	var runs, wickets int
	err := tx.QueryRowx(`
		UPDATE innings
		SET total_runs = COALESCE(total_runs, 0) + $2,
			total_wickets = COALESCE(total_wickets, 0) + $3
		WHERE id = $1
		RETURNING total_runs, total_wickets
	`, inningsID, runsDelta, wicketsDelta).Scan(&runs, &wickets)
	if err != nil {
		return 0, 0, err
	}
	return runs, wickets, nil
}

func UpsertBattingStatsTx(tx *sqlx.Tx, matchID, matchPlayerID string, runs int, legalBall, isFour, isSix, isOut bool) error {
	balls, fours, sixes := 0, 0, 0
	if legalBall {
		balls = 1
	}
	if isFour {
		fours = 1
	}
	if isSix {
		sixes = 1
	}
	res, err := tx.Exec(`
		UPDATE player_match_stats
		SET runs_scored = COALESCE(runs_scored, 0) + $3,
			balls_faced = COALESCE(balls_faced, 0) + $4,
			fours = COALESCE(fours, 0) + $5,
			sixes = COALESCE(sixes, 0) + $6,
			is_out = COALESCE(is_out, FALSE) OR $7,
			updated_at = NOW()
		WHERE match_id = $1 AND team_player_id = $2
	`, matchID, matchPlayerID, runs, balls, fours, sixes, isOut)
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
		INSERT INTO player_match_stats (match_id, player_id, team_player_id, runs_scored, balls_faced, fours, sixes, is_out, updated_at)
		SELECT $1, tp.player_id, tp.id, $3, $4, $5, $6, $7, NOW()
		FROM team_players tp
		WHERE tp.id = $2 AND tp.player_id IS NOT NULL
	`, matchID, matchPlayerID, runs, balls, fours, sixes, isOut)
	return err
}

func UpsertBowlingStatsTx(tx *sqlx.Tx, matchID, matchPlayerID string, runsConceded int, legalBall, isWicket bool) error {
	legalDelta, wickets := 0, 0
	if legalBall {
		legalDelta = 1
	}
	if isWicket {
		wickets = 1
	}
	res, err := tx.Exec(`
		UPDATE player_match_stats
		SET runs_conceded = COALESCE(runs_conceded, 0) + $3,
			wickets_taken = COALESCE(wickets_taken, 0) + $4,
			legal_balls_bowled = COALESCE(legal_balls_bowled, 0) + $5::INTEGER,
			overs_bowled = FLOOR((COALESCE(legal_balls_bowled, 0) + $5::INTEGER) / 6.0)
				+ MOD(COALESCE(legal_balls_bowled, 0) + $5::INTEGER, 6)::NUMERIC / 10.0,
			updated_at = NOW()
		WHERE match_id = $1 AND team_player_id = $2
	`, matchID, matchPlayerID, runsConceded, wickets, legalDelta)
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
		INSERT INTO player_match_stats (match_id, player_id, team_player_id, runs_conceded, wickets_taken, legal_balls_bowled, overs_bowled, updated_at)
		SELECT $1, tp.player_id, tp.id, $3, $4, $5::INTEGER, FLOOR($5::INTEGER / 6.0) + MOD($5::INTEGER, 6)::NUMERIC / 10.0, NOW()
		FROM team_players tp
		WHERE tp.id = $2 AND tp.player_id IS NOT NULL
	`, matchID, matchPlayerID, runsConceded, wickets, legalDelta)
	return err
}

func UpsertFantasyPointsTx(tx *sqlx.Tx, matchID, matchPlayerID string, points int, bucket string) error {
	if bucket != "batting_points" && bucket != "bowling_points" && bucket != "fielding_points" {
		return nil
	}
	query := `UPDATE player_match_stats
		SET ` + bucket + ` = COALESCE(` + bucket + `, 0) + $3,
			fantasy_points = COALESCE(fantasy_points, 0) + $3,
			updated_at = NOW()
		WHERE match_id = $1 AND team_player_id = $2`
	res, err := tx.Exec(query, matchID, matchPlayerID, points)
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
	_, err = tx.Exec(`INSERT INTO player_match_stats (match_id, player_id, team_player_id, `+bucket+`, fantasy_points, updated_at)
		SELECT $1, tp.player_id, tp.id, $3, $3, NOW()
		FROM team_players tp
		WHERE tp.id = $2 AND tp.player_id IS NOT NULL`, matchID, matchPlayerID, points)
	return err
}

func InsertPointEventTx(tx *sqlx.Tx, matchID, matchPlayerID, ballEventID, category, ruleName string, points int) error {
	_, err := tx.Exec(`INSERT INTO point_events (match_id, user_id, ball_event_id, category, rule_name, points)
		SELECT $1, tp.player_id, $2, $3::point_category, $4, $5
		FROM team_players tp
		WHERE tp.id = $6 AND tp.player_id IS NOT NULL`,
		matchID, ballEventID, category, ruleName, points, matchPlayerID)
	return err
}

func ApplyResultPointsTx(tx *sqlx.Tx, matchID, winnerMatchTeamID string) error {
	_, err := tx.Exec(`
		UPDATE player_match_stats pms
		SET result_points = COALESCE(result_points, 0) + 5,
			fantasy_points = COALESCE(fantasy_points, 0) + 5,
			updated_at = NOW()
		FROM team_players tp
		WHERE pms.match_id = $1 AND pms.team_player_id = tp.id
		  AND tp.team_id = $2 AND tp.is_playing_xi = TRUE AND tp.deleted_at IS NULL
	`, matchID, winnerMatchTeamID)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`
		UPDATE player_match_stats pms
		SET result_points = COALESCE(result_points, 0) - 5,
			fantasy_points = COALESCE(fantasy_points, 0) - 5,
			updated_at = NOW()
		FROM team_players tp
		WHERE pms.match_id = $1 AND pms.team_player_id = tp.id
		  AND tp.team_id <> $2 AND tp.is_playing_xi = TRUE AND tp.deleted_at IS NULL
	`, matchID, winnerMatchTeamID)
	return err
}

func ApplyNotOutBonus(matchID, inningsID string) error {
	_, err := database.DB.Exec(`
		UPDATE player_match_stats pms
		SET batting_points = COALESCE(batting_points, 0) + 5,
			fantasy_points = COALESCE(fantasy_points, 0) + 5,
			updated_at = NOW()
		WHERE pms.match_id = $1
		  AND pms.team_player_id IN (
			  SELECT DISTINCT be.striker_id
			  FROM ball_events be
			  WHERE be.match_id = $1 AND be.innings_id = $2
			    AND be.striker_id IS NOT NULL AND be.is_deleted = FALSE
		  )
		  AND COALESCE(pms.is_out, FALSE) = FALSE
	`, matchID, inningsID)
	return err
}

func GetMatchStatus(matchID string) (string, error) {
	var status string
	err := database.DB.Get(&status, `SELECT status FROM matches WHERE id = $1`, matchID)
	return status, err
}

func GetFirstInningsRuns(matchID string) (int, error) {
	var runs int
	err := database.DB.Get(&runs, `SELECT COALESCE(total_runs, 0) FROM innings WHERE match_id = $1 AND innings_no = 1 LIMIT 1`, matchID)
	return runs, err
}

func GetInningsRunsByNo(matchID string, inningsNo int) (int, error) {
	var runs int
	err := database.DB.Get(&runs, `SELECT COALESCE(total_runs, 0) FROM innings WHERE match_id = $1 AND innings_no = $2 LIMIT 1`, matchID, inningsNo)
	return runs, err
}

func IsNoRows(err error) bool { return err == sql.ErrNoRows }
