package repositories

import (
	"CricketDuniya-Backend/internal/dto"
	dbsql "database/sql"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type InningsMeta struct {
	MatchID         string `db:"match_id"`
	InningsNo       int    `db:"innings_no"`
	BattingTeamID   string `db:"batting_team_id"`
	BowlingTeamID   string `db:"bowling_team_id"`
	OversPerInnings int    `db:"overs_per_innings"`
	TargetRuns      *int   `db:"target_runs"`
	IsSuperOver     bool   `db:"is_super_over"`
	SuperOverNo     *int   `db:"super_over_no"`
}

type UndoBall struct {
	ID         string  `db:"id"`
	StrikerID  *string `db:"striker_id"`
	NonStriker *string `db:"non_striker_id"`
	BowlerID   *string `db:"bowler_id"`
	IsFreeHit  bool    `db:"is_free_hit"`
}

type RebuiltInningsTotals struct {
	Runs    int `db:"runs"`
	Wickets int `db:"wickets"`
	Legal   int `db:"legal"`
}
type MatchTeams struct {
	TeamAID string `db:"team_a_id"`
	TeamBID string `db:"team_b_id"`
}

func GetInningsMetaTx(tx *sqlx.Tx, inningsID uuid.UUID) (*InningsMeta, error) {
	var m InningsMeta
	sql := `
		SELECT i.match_id, i.innings_no, i.batting_team_id, i.bowling_team_id,
			   CASE WHEN so.id IS NULL THEN m.overs_per_innings ELSE 1 END AS overs_per_innings,
			   i.target_runs,
			   (so.id IS NOT NULL) AS is_super_over,
			   so.super_over_no
		FROM innings i
		JOIN matches m ON m.id = i.match_id
		LEFT JOIN super_overs so ON so.innings_id = i.id
		WHERE i.id = $1
	`
	err := tx.Get(&m, sql, inningsID)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func CountSecondInningsTx(tx *sqlx.Tx, matchID string) (int, error) {
	var c int
	sql := `SELECT COUNT(1) FROM innings WHERE match_id = $1 AND innings_no = 2`
	err := tx.Get(&c, sql, matchID)
	return c, err
}

func InsertSecondInningsTx(tx *sqlx.Tx, matchID, battingTeamID, bowlingTeamID string, targetRuns int) error {
	sql := `
		INSERT INTO innings (
			match_id, batting_team_id, bowling_team_id, innings_no, target_runs, status, started_at
		) VALUES ($1, $2, $3, 2, $4, 'live', NOW())
	`
	_, err := tx.Exec(sql, matchID, battingTeamID, bowlingTeamID, targetRuns)
	return err
}

func CountInningsByNoTx(tx *sqlx.Tx, matchID string, inningsNo int) (int, error) {
	var c int
	sql := `SELECT COUNT(1) FROM innings WHERE match_id = $1 AND innings_no = $2`
	err := tx.Get(&c, sql, matchID, inningsNo)
	return c, err
}

func InsertInningsTx(tx *sqlx.Tx, matchID, battingTeamID, bowlingTeamID string, inningsNo int, targetRuns *int) error {
	sql := `
		INSERT INTO innings (
			match_id, batting_team_id, bowling_team_id, innings_no, target_runs, status, started_at
		) VALUES ($1, $2, $3, $4, $5, 'live', NOW())
	`
	_, err := tx.Exec(sql, matchID, battingTeamID, bowlingTeamID, inningsNo, targetRuns)
	return err
}

func GetMatchTeamsTx(tx *sqlx.Tx, matchID string) (*MatchTeams, error) {
	var t MatchTeams
	sql := `SELECT team_a_id, team_b_id FROM matches WHERE id = $1`
	if err := tx.Get(&t, sql, matchID); err != nil {
		return nil, err
	}
	return &t, nil
}

//func GetMaxSuperOverNoTx(tx *sqlx.Tx, matchID string) (int, error) {
//	var n int
//	sql := `SELECT COALESCE(MAX(super_over_no), 0) FROM super_overs WHERE match_id = $1`
//	err := tx.Get(&n, sql, matchID)
//	return n, err
//}

func LinkSuperOverToInningsTx(tx *sqlx.Tx, matchID string, inningsID uuid.UUID, battingTeamID, bowlingTeamID string, superOverNo int) error {
	sql := `
		INSERT INTO super_overs (
			match_id, innings_id, batting_match_team_id, bowling_match_team_id, super_over_no, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
	`
	_, err := tx.Exec(sql, matchID, inningsID, battingTeamID, bowlingTeamID, superOverNo)
	return err
}

func GetInningsIDByNoTx(tx *sqlx.Tx, matchID string, inningsNo int) (uuid.UUID, error) {
	var id uuid.UUID
	sql := `SELECT id FROM innings WHERE match_id = $1 AND innings_no = $2 LIMIT 1`
	err := tx.Get(&id, sql, matchID, inningsNo)
	return id, err
}

//func GetInningsRunsByNoTx(tx *sqlx.Tx, matchID string, inningsNo int) (int, error) {
//	var runs int
//	sql := `SELECT COALESCE(total_runs, 0) FROM innings WHERE match_id = $1 AND innings_no = $2 LIMIT 1`
//	err := tx.Get(&runs, sql, matchID, inningsNo)
//	return runs, err
//}

func GetPreviousOverBowlerIDTx(tx *sqlx.Tx, inningsID uuid.UUID) (string, error) {
	var bowlerID dbsql.NullString
	sql := `
		SELECT bowler_id
		FROM ball_events
		WHERE innings_id = $1 AND is_deleted = FALSE
		ORDER BY created_at DESC
		LIMIT 1
	`
	err := tx.Get(&bowlerID, sql, inningsID)
	if err != nil {
		if err == dbsql.ErrNoRows {
			return "", nil
		}
		return "", err
	}
	if !bowlerID.Valid {
		return "", nil
	}
	return bowlerID.String, nil
}

func GetInningsStateForUpdateTx(tx *sqlx.Tx, inningsID uuid.UUID, dst interface{}) error {
	sql := `
		SELECT
			innings_id,
			striker_id,
			non_striker_id,
			bowler_id,
			is_free_hit,
			total_runs,
			total_wickets,
			legal_balls,
			current_over,
			current_ball,
			status,
			updated_at
		FROM innings_state
		WHERE innings_id = $1
		FOR UPDATE
	`
	return tx.Get(dst, sql, inningsID)
}

func InsertInitialInningsStateTx(tx *sqlx.Tx, req dto.BallInputRequest) error {
	sql := `
		INSERT INTO innings_state (innings_id, striker_id, non_striker_id, bowler_id, total_runs, total_wickets, legal_balls, current_over, current_ball, is_free_hit, status)
		VALUES ($1,$2,$3,$4,0,0,0,0,0,FALSE,'live')
	`
	_, err := tx.Exec(sql, req.InningsID, req.StrikerID, req.NonStrikerID, req.BowlerID)
	return err
}

func CountStrikersOnTeamTx(tx *sqlx.Tx, strikerID, nonStrikerID, teamID string) (int, error) {
	var c int
	sql := `SELECT COUNT(1) FROM team_players WHERE id IN ($1,$2) AND team_id = $3 AND deleted_at IS NULL`
	err := tx.Get(&c, sql, strikerID, nonStrikerID, teamID)
	return c, err
}

func CountActiveTeamPlayersTx(tx *sqlx.Tx, teamID string) (int, error) {
	var c int
	sql := `SELECT COUNT(1) FROM team_players WHERE team_id = $1 AND deleted_at IS NULL`
	err := tx.Get(&c, sql, teamID)
	return c, err
}

func CountPlayerOnTeamTx(tx *sqlx.Tx, playerID, teamID string) (int, error) {
	var c int
	sql := `SELECT COUNT(1) FROM team_players WHERE id = $1 AND team_id = $2 AND deleted_at IS NULL`
	err := tx.Get(&c, sql, playerID, teamID)
	return c, err
}

func CountRetiredHurtBattersTx(tx *sqlx.Tx, inningsID uuid.UUID) (int, error) {
	var c int
	sql := `
		SELECT COUNT(DISTINCT dismissed_player_id)
		FROM ball_events
		WHERE innings_id = $1
		  AND is_deleted = FALSE
		  AND dismissed_player_id IS NOT NULL
		  AND (
			  LOWER(COALESCE(ball_type, '')) = 'retired_hurt'
			  OR LOWER(COALESCE(dismissal_type, '')) = 'retired_hurt'
		  )
	`
	err := tx.Get(&c, sql, inningsID)
	return c, err
}

func CountBowlerOnTeamTx(tx *sqlx.Tx, bowlerID, teamID string) (int, error) {
	var c int
	sql := `SELECT COUNT(1) FROM team_players WHERE id = $1 AND team_id = $2 AND deleted_at IS NULL`
	err := tx.Get(&c, sql, bowlerID, teamID)
	return c, err
}

func SaveBallEventTx(tx *sqlx.Tx, req dto.BallRequest) (string, error) {
	var id string
	sql := `
		INSERT INTO ball_events (
			id, innings_id, match_id, striker_id, non_striker_id, bowler_id,
			ball_no, delivery_no, ball_type, runs_scored, runs_off_bat, extras, total_runs,
			is_dot_ball, is_boundary_four, is_boundary_six, is_wicket, dismissal_type,
			dismissed_player_id, fielder_id, wides, no_balls, byes, leg_byes, is_free_hit, created_at
		) VALUES (
			gen_random_uuid(), $1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,
			$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24, NOW()
		) RETURNING id
	`
	err := tx.Get(&id, sql, req.InningsID, req.MatchID, req.StrikerID, req.NonStrikerID, req.BowlerID,
		req.BallNo, req.DeliveryNo, req.BallType, req.RunsScored, req.RunsOffBat, req.Extras, req.TotalRuns,
		req.IsDotBall, req.IsBoundaryFour, req.IsBoundarySix, req.IsWicket, req.DismissalType, req.DismissedPlayerID, req.FielderID,
		req.Wides, req.NoBalls, req.Byes, req.LegByes, req.IsFreeHit)
	if err != nil {
		return "", err
	}
	return id, nil
}

func UpdateInningsStateAfterBallTx(tx *sqlx.Tx, inningsID uuid.UUID, strikerID, nonStrikerID, bowlerID string, runs, wickets, legalBalls, over, ball int, isFreeHit bool, status string) error {
	sql := `
		UPDATE innings_state
		SET striker_id = $2, non_striker_id = $3, bowler_id = $4,
			total_runs = $5, total_wickets = $6, legal_balls = $7,
			current_over = $8, current_ball = $9, is_free_hit = $10, status = $11, updated_at = NOW()
		WHERE innings_id = $1
	`
	_, err := tx.Exec(sql, inningsID, strikerID, nonStrikerID, bowlerID, runs, wickets, legalBalls, over, ball, isFreeHit, status)
	return err
}

func UpdateInningsTotalsStatusTx(tx *sqlx.Tx, inningsID uuid.UUID, runs, wickets int, status string) error {
	sql := `
		UPDATE innings
		SET total_runs = $2, total_wickets = $3, status = $4
		WHERE id = $1
	`
	_, err := tx.Exec(sql, inningsID, runs, wickets, status)
	return err
}

func GetInningsStateTx(tx *sqlx.Tx, inningsID uuid.UUID, dst interface{}) error {
	sql := `
		SELECT
			innings_id,
			striker_id,
			non_striker_id,
			bowler_id,
			is_free_hit,
			total_runs,
			total_wickets,
			legal_balls,
			current_over,
			current_ball,
			status,
			updated_at
		FROM innings_state
		WHERE innings_id = $1
	`
	return tx.Get(dst, sql, inningsID)
}

func GetLastUndoBallTx(tx *sqlx.Tx, inningsID uuid.UUID) (*UndoBall, error) {
	var out UndoBall
	sql := `
		SELECT id, striker_id, non_striker_id, bowler_id, is_free_hit
		FROM ball_events
		WHERE innings_id = $1 AND is_deleted = FALSE
		ORDER BY created_at DESC
		LIMIT 1
	`
	err := tx.Get(&out, sql, inningsID)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func SoftDeleteBallTx(tx *sqlx.Tx, ballID string) error {
	sql := `UPDATE ball_events SET is_deleted = TRUE, deleted_at = NOW() WHERE id = $1`
	_, err := tx.Exec(sql, ballID)
	return err
}

func DeletePointEventsByBallTx(tx *sqlx.Tx, ballID string) error {
	sql := `DELETE FROM point_events WHERE ball_event_id = $1`
	_, err := tx.Exec(sql, ballID)
	return err
}

func GetRebuiltInningsTotalsTx(tx *sqlx.Tx, inningsID uuid.UUID) (*RebuiltInningsTotals, error) {
	var rebuilt RebuiltInningsTotals
	sql := `
		SELECT COALESCE(SUM(total_runs),0) AS runs,
			   COALESCE(SUM(CASE WHEN is_wicket AND LOWER(COALESCE(dismissal_type, '')) <> 'retired_hurt' AND LOWER(COALESCE(ball_type, '')) <> 'retired_hurt' THEN 1 ELSE 0 END),0) AS wickets,
			   COALESCE(SUM(CASE WHEN COALESCE(ball_type, 'normal') NOT IN ('wide','no_ball','dead_ball','retired_hurt') THEN 1 ELSE 0 END),0) AS legal
		FROM ball_events
		WHERE innings_id = $1 AND is_deleted = FALSE
	`
	err := tx.Get(&rebuilt, sql, inningsID)
	if err != nil {
		return nil, err
	}
	return &rebuilt, nil
}

func CountFutureBallsTx(tx *sqlx.Tx, matchID string, inningsNo int) (int, error) {
	var c int
	sql := `
		SELECT COUNT(1)
		FROM innings i
		JOIN ball_events be ON be.innings_id = i.id
		WHERE i.match_id = $1 AND i.innings_no > $2 AND be.is_deleted = FALSE
	`
	err := tx.Get(&c, sql, matchID, inningsNo)
	return c, err
}

func UpdateInningsStateAfterUndoTx(tx *sqlx.Tx, inningsID uuid.UUID, strikerID, nonStrikerID, bowlerID *string, runs, wickets, legalBalls, over, ball int, isFreeHit bool) error {
	sql := `
		UPDATE innings_state
		SET striker_id = $2, non_striker_id = $3, bowler_id = $4,
			total_runs = $5, total_wickets = $6, legal_balls = $7,
			current_over = $8, current_ball = $9, is_free_hit = $10, status = 'live', updated_at = NOW()
		WHERE innings_id = $1
	`
	_, err := tx.Exec(sql, inningsID, strikerID, nonStrikerID, bowlerID, runs, wickets, legalBalls, over, ball, isFreeHit)
	return err
}

func ReopenInningsTx(tx *sqlx.Tx, inningsID uuid.UUID, runs, wickets int) error {
	sql := `
		UPDATE innings
		SET total_runs = $2, total_wickets = $3, status = 'live', completed_at = NULL
		WHERE id = $1
	`
	_, err := tx.Exec(sql, inningsID, runs, wickets)
	return err
}

func DeleteFutureInningsStateTx(tx *sqlx.Tx, matchID string, inningsNo int) error {
	sql := `
		DELETE FROM innings_state
		WHERE innings_id IN (
			SELECT id FROM innings WHERE match_id = $1 AND innings_no > $2
		)
	`
	_, err := tx.Exec(sql, matchID, inningsNo)
	return err
}

func DeleteFutureInningsTx(tx *sqlx.Tx, matchID string, inningsNo int) error {
	sql := `DELETE FROM innings WHERE match_id = $1 AND innings_no > $2`
	_, err := tx.Exec(sql, matchID, inningsNo)
	return err
}

func ReopenMatchAfterUndoTx(tx *sqlx.Tx, matchID string) error {
	sql := `
		UPDATE player_match_stats
		SET fantasy_points = COALESCE(fantasy_points, 0) - COALESCE(result_points, 0),
			result_points = 0,
			updated_at = NOW()
		WHERE match_id = $1
	`
	if _, err := tx.Exec(sql, matchID); err != nil {
		return err
	}
	sql = `
		UPDATE matches
		SET status = 'live',
			winner_match_team_id = NULL,
			completed_at = NULL,
			player_of_match_user_id = NULL,
			worst_player_user_id = NULL
		WHERE id = $1
	`
	_, err := tx.Exec(sql, matchID)
	return err
}

func RebuildPlayerMatchStatsTx(tx *sqlx.Tx, matchID string) error {
	sql := `
		INSERT INTO player_match_stats (match_id, player_id, team_player_id, updated_at)
		SELECT $1, tp.player_id, tp.id, NOW()
		FROM team_players tp
		WHERE tp.player_id IS NOT NULL
		  AND tp.id IN (
			SELECT be.striker_id FROM ball_events be WHERE be.match_id = $1 AND be.is_deleted = FALSE AND be.striker_id IS NOT NULL
			UNION
			SELECT be.non_striker_id FROM ball_events be WHERE be.match_id = $1 AND be.is_deleted = FALSE AND be.non_striker_id IS NOT NULL
			UNION
			SELECT be.bowler_id FROM ball_events be WHERE be.match_id = $1 AND be.is_deleted = FALSE AND be.bowler_id IS NOT NULL
			UNION
			SELECT be.dismissed_player_id FROM ball_events be WHERE be.match_id = $1 AND be.is_deleted = FALSE AND be.dismissed_player_id IS NOT NULL
			UNION
			SELECT tp2.id
			FROM point_events pe
			JOIN team_players tp2 ON tp2.player_id = pe.user_id
			JOIN matches m ON m.id = pe.match_id
			WHERE pe.match_id = $1 AND (tp2.team_id = m.team_a_id OR tp2.team_id = m.team_b_id)
		  )
		  AND NOT EXISTS (
			SELECT 1 FROM player_match_stats pms WHERE pms.match_id = $1 AND pms.team_player_id = tp.id
		  )
	`
	_, err := tx.Exec(sql, matchID)
	if err != nil {
		return err
	}

	sql = `
		UPDATE player_match_stats
		SET runs_scored = 0, balls_faced = 0, fours = 0, sixes = 0, strike_rate = 0, is_out = FALSE,
			overs_bowled = 0, legal_balls_bowled = 0, maidens = 0, runs_conceded = 0, wickets_taken = 0, economy = 0,
			catches = 0, stumping = 0, runouts = 0,
			batting_points = 0, bowling_points = 0, fielding_points = 0, fantasy_points = COALESCE(result_points, 0),
			updated_at = NOW()
		WHERE match_id = $1
	`
	_, err = tx.Exec(sql, matchID)
	if err != nil {
		return err
	}

	sql = `
		WITH batting AS (
			SELECT be.striker_id AS team_player_id,
				COALESCE(SUM(GREATEST(COALESCE(be.total_runs, 0) - COALESCE(be.extras, 0), 0)), 0)::INT AS runs_scored,
				COALESCE(SUM(CASE WHEN COALESCE(be.ball_type, 'normal') NOT IN ('wide', 'no_ball', 'dead_ball', 'retired_hurt') THEN 1 ELSE 0 END), 0)::INT AS balls_faced,
				COALESCE(SUM(CASE WHEN GREATEST(COALESCE(be.total_runs, 0) - COALESCE(be.extras, 0), 0) = 4 THEN 1 ELSE 0 END), 0)::INT AS fours,
				COALESCE(SUM(CASE WHEN GREATEST(COALESCE(be.total_runs, 0) - COALESCE(be.extras, 0), 0) = 6 THEN 1 ELSE 0 END), 0)::INT AS sixes
			FROM ball_events be
			WHERE be.match_id = $1 AND be.is_deleted = FALSE AND be.striker_id IS NOT NULL
			GROUP BY be.striker_id
		)
		UPDATE player_match_stats pms
		SET runs_scored = batting.runs_scored, balls_faced = batting.balls_faced, fours = batting.fours, sixes = batting.sixes,
			strike_rate = CASE WHEN batting.balls_faced > 0 THEN ROUND(batting.runs_scored::NUMERIC * 100 / batting.balls_faced, 2) ELSE 0 END,
			updated_at = NOW()
		FROM batting
		WHERE pms.match_id = $1 AND pms.team_player_id = batting.team_player_id
	`
	_, err = tx.Exec(sql, matchID)
	if err != nil {
		return err
	}

	sql = `
		WITH dismissed AS (
			SELECT DISTINCT dismissed_player_id AS team_player_id
			FROM ball_events
			WHERE match_id = $1
			  AND is_deleted = FALSE
			  AND is_wicket = TRUE
			  AND dismissed_player_id IS NOT NULL
			  AND LOWER(COALESCE(ball_type, '')) <> 'retired_hurt'
			  AND LOWER(COALESCE(dismissal_type, '')) <> 'retired_hurt'
		)
		UPDATE player_match_stats pms
		SET is_out = TRUE, updated_at = NOW()
		FROM dismissed
		WHERE pms.match_id = $1 AND pms.team_player_id = dismissed.team_player_id
	`
	_, err = tx.Exec(sql, matchID)
	if err != nil {
		return err
	}

	sql = `
		WITH bowling AS (
			SELECT be.bowler_id AS team_player_id,
				COALESCE(SUM(GREATEST(COALESCE(be.total_runs, 0) - COALESCE(be.byes, 0) - COALESCE(be.leg_byes, 0), 0)), 0)::INT AS runs_conceded,
				COALESCE(SUM(CASE WHEN COALESCE(be.is_wicket, FALSE) AND LOWER(COALESCE(be.ball_type, '')) <> 'retired_hurt' AND LOWER(COALESCE(be.dismissal_type, '')) <> 'retired_hurt' THEN 1 ELSE 0 END), 0)::INT AS wickets_taken,
				COALESCE(SUM(CASE WHEN COALESCE(be.ball_type, 'normal') NOT IN ('wide', 'no_ball', 'dead_ball', 'retired_hurt') THEN 1 ELSE 0 END), 0)::INT AS legal_balls_bowled
			FROM ball_events be
			WHERE be.match_id = $1 AND be.is_deleted = FALSE AND be.bowler_id IS NOT NULL
			GROUP BY be.bowler_id
		), maiden_overs AS (
			SELECT over_summary.team_player_id,
				COUNT(*) FILTER (WHERE over_summary.legal_balls = 6 AND over_summary.runs_conceded = 0)::INT AS maidens
			FROM (
				SELECT be.bowler_id AS team_player_id, be.innings_id, be.ball_no,
					SUM(CASE WHEN COALESCE(be.ball_type, 'normal') NOT IN ('wide', 'no_ball', 'dead_ball', 'retired_hurt') THEN 1 ELSE 0 END) AS legal_balls,
					SUM(GREATEST(COALESCE(be.total_runs, 0) - COALESCE(be.byes, 0) - COALESCE(be.leg_byes, 0), 0)) AS runs_conceded
				FROM ball_events be
				WHERE be.match_id = $1 AND be.is_deleted = FALSE AND be.bowler_id IS NOT NULL
				GROUP BY be.bowler_id, be.innings_id, be.ball_no
			) over_summary
			GROUP BY over_summary.team_player_id
		)
		UPDATE player_match_stats pms
		SET runs_conceded = bowling.runs_conceded, wickets_taken = bowling.wickets_taken,
			legal_balls_bowled = bowling.legal_balls_bowled,
			overs_bowled = FLOOR(bowling.legal_balls_bowled / 6.0) + MOD(bowling.legal_balls_bowled, 6)::NUMERIC / 10.0,
			maidens = COALESCE(maiden_overs.maidens, 0),
			economy = CASE WHEN bowling.legal_balls_bowled > 0 THEN ROUND(bowling.runs_conceded::NUMERIC * 6 / bowling.legal_balls_bowled, 2) ELSE 0 END,
			updated_at = NOW()
		FROM bowling
		LEFT JOIN maiden_overs ON maiden_overs.team_player_id = bowling.team_player_id
		WHERE pms.match_id = $1 AND pms.team_player_id = bowling.team_player_id
	`
	_, err = tx.Exec(sql, matchID)
	if err != nil {
		return err
	}

	sql = `
		WITH fielding AS (
			SELECT fielder_id AS team_player_id,
				COALESCE(SUM(CASE WHEN dismissal_type = 'caught' THEN 1 ELSE 0 END), 0)::INT AS catches,
				COALESCE(SUM(CASE WHEN dismissal_type = 'stumped' THEN 1 ELSE 0 END), 0)::INT AS stumping,
				COALESCE(SUM(CASE WHEN dismissal_type = 'run_out' THEN 1 ELSE 0 END), 0)::INT AS runouts
			FROM ball_events
			WHERE match_id = $1 AND is_deleted = FALSE AND fielder_id IS NOT NULL
			GROUP BY fielder_id
		)
		UPDATE player_match_stats pms
		SET catches = fielding.catches, stumping = fielding.stumping, runouts = fielding.runouts, updated_at = NOW()
		FROM fielding
		WHERE pms.match_id = $1 AND pms.team_player_id = fielding.team_player_id
	`
	_, err = tx.Exec(sql, matchID)
	if err != nil {
		return err
	}

	sql = `WITH points AS (
			SELECT tp.id AS team_player_id,
				COALESCE(SUM(CASE WHEN pe.category = 'batting'::point_category THEN pe.points ELSE 0 END), 0)::INT AS batting_points,
				COALESCE(SUM(CASE WHEN pe.category = 'bowling'::point_category THEN pe.points ELSE 0 END), 0)::INT AS bowling_points,
				COALESCE(SUM(CASE WHEN pe.category = 'fielding'::point_category THEN pe.points ELSE 0 END), 0)::INT AS fielding_points,
				COALESCE(SUM(pe.points), 0)::INT AS fantasy_points
			FROM point_events pe
			JOIN team_players tp ON tp.player_id = pe.user_id
			JOIN matches m ON m.id = pe.match_id
			WHERE pe.match_id = $1 AND (tp.team_id = m.team_a_id OR tp.team_id = m.team_b_id)
			GROUP BY tp.id
		)
		UPDATE player_match_stats pms
		SET batting_points = points.batting_points, bowling_points = points.bowling_points, fielding_points = points.fielding_points,
			fantasy_points = points.fantasy_points + COALESCE(pms.result_points, 0), updated_at = NOW()
		FROM points
		WHERE pms.match_id = $1 AND pms.team_player_id = points.team_player_id
	`
	_, err = tx.Exec(sql, matchID)
	return err
}
