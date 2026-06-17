package repositories

import (
	"CricketDuniya-Backend/internal/database"

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
	sql := `
		SELECT
			COALESCE(SUM(CASE WHEN ball_type IN ('wide','bye','leg_bye') THEN 0 ELSE GREATEST(total_runs - extras, 0) END), 0) AS runs_before,
			COALESCE(SUM(CASE WHEN ball_type NOT IN ('wide','no_ball','dead_ball','retired_hurt') THEN 1 ELSE 0 END), 0) AS balls_before
		FROM ball_events
		WHERE match_id = $1 AND innings_id = $2 AND striker_id = $3 AND is_deleted = FALSE
	`
	err := database.DB.Get(&out, sql, matchID, inningsID, strikerID)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func GetBowlerTotals(matchID, inningsID, bowlerID string) (*BowlerTotals, error) {
	var out BowlerTotals
	sql := `
		SELECT
			COALESCE(SUM(CASE WHEN ball_type NOT IN ('wide','no_ball','dead_ball','retired_hurt') THEN 1 ELSE 0 END), 0) AS legal_balls_before,
			COALESCE(SUM(total_runs - byes - leg_byes), 0) AS runs_conceded_before
		FROM ball_events
		WHERE match_id = $1 AND innings_id = $2 AND bowler_id = $3 AND is_deleted = FALSE
	`
	err := database.DB.Get(&out, sql, matchID, inningsID, bowlerID)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func GetBowlerOverTotals(matchID, inningsID, bowlerID string, ballNo int) (*BowlerOverTotals, error) {
	var out BowlerOverTotals
	sql := `
		SELECT
			COALESCE(SUM(CASE WHEN ball_type NOT IN ('wide','no_ball','dead_ball','retired_hurt') THEN 1 ELSE 0 END), 0) AS over_legal_before,
			COALESCE(SUM(total_runs - byes - leg_byes), 0) AS over_runs_before,
			COALESCE(SUM(CASE WHEN ball_type IN ('wide','bye','leg_bye') THEN 0 WHEN GREATEST(total_runs - extras, 0) IN (4, 6) THEN 1 ELSE 0 END), 0) AS over_boundaries_before
		FROM ball_events
		WHERE match_id = $1 AND innings_id = $2 AND bowler_id = $3 AND ball_no = $4 AND is_deleted = FALSE
	`
	err := database.DB.Get(&out, sql, matchID, inningsID, bowlerID, ballNo)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func ListRecentRunsOffBatForStriker(matchID, inningsID, strikerID string, limit int) ([]int, error) {
	sql := `
		SELECT CASE WHEN ball_type IN ('wide','bye','leg_bye') THEN 0 ELSE GREATEST(total_runs - extras, 0) END AS runs_off_bat
		FROM ball_events
		WHERE match_id = $1 AND innings_id = $2 AND striker_id = $3 AND is_deleted = FALSE
		ORDER BY created_at DESC
		LIMIT $4
	`
	out := make([]int, 0, limit)
	err := database.DB.Select(&out, sql, matchID, inningsID, strikerID, limit)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func ListRecentRunsOffBatForBowler(matchID, inningsID, bowlerID string, limit int) ([]int, error) {
	sql := `
		SELECT CASE WHEN ball_type IN ('wide','bye','leg_bye') THEN 0 ELSE GREATEST(total_runs - extras, 0) END AS runs_off_bat
		FROM ball_events
		WHERE match_id = $1 AND innings_id = $2 AND bowler_id = $3 AND is_deleted = FALSE
		ORDER BY created_at DESC
		LIMIT $4
	`
	out := make([]int, 0, limit)
	err := database.DB.Select(&out, sql, matchID, inningsID, bowlerID, limit)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func ListRecentBallTypesForBowler(matchID, inningsID, bowlerID string, limit int) ([]string, error) {
	sql := `
		SELECT ball_type
		FROM ball_events
		WHERE match_id = $1 AND innings_id = $2 AND bowler_id = $3 AND is_deleted = FALSE
		ORDER BY created_at DESC
		LIMIT $4
	`
	out := make([]string, 0, limit)
	err := database.DB.Select(&out, sql, matchID, inningsID, bowlerID, limit)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func ListRecentWicketsForBowlerInOver(matchID, inningsID, bowlerID string, ballNo, limit int) ([]bool, error) {
	sql := `
		SELECT is_wicket
		FROM ball_events
		WHERE match_id = $1 AND innings_id = $2 AND bowler_id = $3 AND ball_no = $4 AND is_deleted = FALSE
		ORDER BY created_at DESC
		LIMIT $5
	`
	out := make([]bool, 0, limit)
	err := database.DB.Select(&out, sql, matchID, inningsID, bowlerID, ballNo, limit)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func UpdateInningsTotalsTx(tx *sqlx.Tx, inningsID string, runsDelta, wicketsDelta int) (int, int, error) {
	var totals struct {
		Runs    int `db:"total_runs"`
		Wickets int `db:"total_wickets"`
	}
	sql := `
		UPDATE innings
		SET total_runs = COALESCE(total_runs, 0) + $2,
			total_wickets = COALESCE(total_wickets, 0) + $3
		WHERE id = $1
		RETURNING total_runs, total_wickets
	`
	err := tx.Get(&totals, sql, inningsID, runsDelta, wicketsDelta)
	if err != nil {
		return 0, 0, err
	}
	return totals.Runs, totals.Wickets, nil
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
	sql := `
		UPDATE player_match_stats
		SET runs_scored = COALESCE(runs_scored, 0) + $3,
			balls_faced = COALESCE(balls_faced, 0) + $4,
			fours = COALESCE(fours, 0) + $5,
			sixes = COALESCE(sixes, 0) + $6,
			is_out = COALESCE(is_out, FALSE) OR $7,
			updated_at = NOW()
		WHERE match_id = $1 AND team_player_id = $2
	`
	res, err := tx.Exec(sql, matchID, matchPlayerID, runs, balls, fours, sixes, isOut)
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
	sql = `
		INSERT INTO player_match_stats (match_id, player_id, team_player_id, runs_scored, balls_faced, fours, sixes, is_out, updated_at)
		SELECT $1, tp.player_id, tp.id, $3, $4, $5, $6, $7, NOW()
		FROM team_players tp
		WHERE tp.id = $2 AND tp.player_id IS NOT NULL
	`
	_, err = tx.Exec(sql, matchID, matchPlayerID, runs, balls, fours, sixes, isOut)
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
	sql := `
UPDATE player_match_stats
SET runs_conceded      = COALESCE(runs_conceded, 0) + $3,
    wickets_taken      = COALESCE(wickets_taken, 0) + $4,
    legal_balls_bowled = COALESCE(legal_balls_bowled, 0) + $5::INTEGER,
    overs_bowled       = FLOOR((COALESCE(legal_balls_bowled, 0) + $5::INTEGER) / 6.0)
        + MOD(COALESCE(legal_balls_bowled, 0) + $5::INTEGER, 6)::NUMERIC / 10.0,
    updated_at         = NOW()
WHERE match_id = $1
  AND team_player_id = $2
    `
	res, err := tx.Exec(sql, matchID, matchPlayerID, runsConceded, wickets, legalDelta)
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
	sql = `
		INSERT INTO player_match_stats (match_id, player_id, team_player_id, runs_conceded, wickets_taken, legal_balls_bowled, overs_bowled, updated_at)
		SELECT $1, tp.player_id, tp.id, $3, $4, $5::INTEGER, FLOOR($5::INTEGER / 6.0) + MOD($5::INTEGER, 6)::NUMERIC / 10.0, NOW()
		FROM team_players tp
		WHERE tp.id = $2 AND tp.player_id IS NOT NULL
	`
	_, err = tx.Exec(sql, matchID, matchPlayerID, runsConceded, wickets, legalDelta)
	return err
}

func UpsertFieldingStatsTx(tx *sqlx.Tx, matchID, matchPlayerID string, isCatch, isStumping, isRunOut bool) error {

	catches, stumpings, runouts := 0, 0, 0

	if isCatch {
		catches = 1
	}
	if isStumping {
		stumpings = 1
	}
	if isRunOut {
		runouts = 1
	}
	if catches == 0 && stumpings == 0 && runouts == 0 {
		return nil
	}

	sql := `
		UPDATE player_match_stats
		SET catches = COALESCE(catches, 0) + $3,
			stumping = COALESCE(stumping, 0) + $4,
			runouts = COALESCE(runouts, 0) + $5,
			updated_at = NOW()
		WHERE match_id = $1 AND team_player_id = $2
	`
	res, err := tx.Exec(sql, matchID, matchPlayerID, catches, stumpings, runouts)
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
	sql = `
		INSERT INTO player_match_stats (match_id, player_id, team_player_id, catches, stumping, runouts, updated_at)
		SELECT $1, tp.player_id, tp.id, $3, $4, $5, NOW()
		FROM team_players tp
		WHERE tp.id = $2 AND tp.player_id IS NOT NULL
	`
	_, err = tx.Exec(sql, matchID, matchPlayerID, catches, stumpings, runouts)
	return err
}

func UpsertFantasyPointsTx(tx *sqlx.Tx, matchID, matchPlayerID string, points int, bucket string) error {
	if bucket != "batting_points" && bucket != "bowling_points" && bucket != "fielding_points" {
		return nil
	}
	sql := `UPDATE player_match_stats
		SET` + bucket + ` = COALESCE(` + bucket + `, 0) + $3,
			fantasy_points = COALESCE(fantasy_points, 0) + $3,
			updated_at = NOW()
		WHERE match_id = $1 AND team_player_id = $2`
	res, err := tx.Exec(sql, matchID, matchPlayerID, points)
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
	sql = `INSERT INTO player_match_stats (match_id, player_id, team_player_id, ` + bucket + `, fantasy_points, updated_at)
		SELECT $1, tp.player_id, tp.id, $3, $3, NOW()
		FROM team_players tp
		WHERE tp.id = $2 AND tp.player_id IS NOT NULL`
	_, err = tx.Exec(sql, matchID, matchPlayerID, points)
	return err
}

func InsertPointEventTx(tx *sqlx.Tx, matchID, matchPlayerID, ballEventID, category, ruleName string, points int) error {
	sql := `INSERT INTO point_events (match_id, user_id, ball_event_id, category, rule_name, points)
		SELECT $1, tp.player_id, $2, $3::point_category, $4, $5
		FROM team_players tp
		WHERE tp.id = $6 AND tp.player_id IS NOT NULL`
	_, err := tx.Exec(sql, matchID, ballEventID, category, ruleName, points, matchPlayerID)
	return err
}

func ApplyResultPointsTx(tx *sqlx.Tx, matchID, winnerMatchTeamID string) error {
	sql := `UPDATE player_match_stats pms
		SET result_points = COALESCE(result_points, 0) + 5,
			fantasy_points = COALESCE(fantasy_points, 0) + 5,
			updated_at = NOW()
		FROM team_players tp
		WHERE pms.match_id = $1 AND pms.team_player_id = tp.id
		  AND tp.team_id = $2 AND tp.is_playing_xi = TRUE AND tp.deleted_at IS NULL
	`
	_, err := tx.Exec(sql, matchID, winnerMatchTeamID)
	if err != nil {
		return err
	}
	sql = `UPDATE player_match_stats pms
		SET result_points = COALESCE(result_points, 0) - 5,
			fantasy_points = COALESCE(fantasy_points, 0) - 5,
			updated_at = NOW()
		FROM team_players tp
		WHERE pms.match_id = $1 AND pms.team_player_id = tp.id
		  AND tp.team_id <> $2 AND tp.is_playing_xi = TRUE AND tp.deleted_at IS NULL
	`
	_, err = tx.Exec(sql, matchID, winnerMatchTeamID)
	return err
}

func ApplyNotOutBonus(matchID, inningsID string) error {
	sql := `UPDATE player_match_stats pms
		SET batting_points = COALESCE(batting_points, 0) + 5,
			fantasy_points = COALESCE(fantasy_points, 0) + 5,
			updated_at = NOW()
		WHERE pms.match_id = $1
		  AND pms.team_player_id IN (
			  SELECT DISTINCT be.striker_id
			  FROM ball_events be
			  WHERE be.match_id = $1 AND be.innings_id = $2
			    AND be.striker_id IS NOT NULL AND be.is_deleted = FALSE
			  UNION
			  SELECT DISTINCT be.non_striker_id
			  FROM ball_events be
			  WHERE be.match_id = $1 AND be.innings_id = $2
			    AND be.non_striker_id IS NOT NULL AND be.is_deleted = FALSE
			  UNION
			  SELECT DISTINCT be.dismissed_player_id
			  FROM ball_events be
			  WHERE be.match_id = $1 AND be.innings_id = $2
			    AND be.dismissed_player_id IS NOT NULL
			    AND be.is_deleted = FALSE
			    AND (
			        LOWER(COALESCE(be.ball_type, '')) = 'retired_hurt'
			        OR LOWER(COALESCE(be.dismissal_type, '')) = 'retired_hurt'
			    )
		  )
		  AND COALESCE(pms.is_out, FALSE) = FALSE
	`
	_, err := database.DB.Exec(sql, matchID, inningsID)
	return err
}

func GetMatchStatus(matchID string) (string, error) {
	var status string
	sql := `SELECT status FROM matches WHERE id = $1`
	err := database.DB.Get(&status, sql, matchID)
	return status, err
}

func GetFirstInningsRuns(matchID string) (int, error) {
	var runs int
	sql := `SELECT COALESCE(total_runs, 0) FROM innings WHERE match_id = $1 AND innings_no = 1 LIMIT 1`
	err := database.DB.Get(&runs, sql, matchID)
	return runs, err
}

func GetInningsRunsByNo(matchID string, inningsNo int) (int, error) {
	var runs int
	sql := `SELECT COALESCE(total_runs, 0) FROM innings WHERE match_id = $1 AND innings_no = $2 LIMIT 1`
	err := database.DB.Get(&runs, sql, matchID, inningsNo)
	return runs, err
}
