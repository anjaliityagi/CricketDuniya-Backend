package repositories

import (
	"CricketDuniya-Backend/internal/database"
	"CricketDuniya-Backend/internal/dto"
	"CricketDuniya-Backend/internal/models"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type MatchTeamSnapshot struct {
	ID           string  `db:"id"`
	MatchID      string  `db:"match_id"`
	SourceTeamID *string `db:"source_team_id"`
	DisplayName  string  `db:"display_name"`
}

func CreateMatch(match *models.Match) error {

	query := `
	INSERT INTO matches (
	           team_a_id,
	                     team_b_id,
		host_user_id,
		location,
		match_date,
		overs_per_innings,
		status
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING id, created_at
	`

	return database.DB.QueryRowx(
		query,
		match.TeamAID,
		match.TeamBID,
		match.HostUserID,
		match.Location,
		match.MatchDate,
		match.OversPerInnings,
		match.Status,
	).Scan(
		&match.ID,
		&match.CreatedAt,
	)
}

func CreateMatchTx(tx *sqlx.Tx, match *models.Match) error {
	query := `
	INSERT INTO matches (
	           team_a_id,
	                     team_b_id,
		host_user_id,
		location,
		match_date,
		overs_per_innings,
		status
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING id, created_at
	`

	return tx.QueryRowx(
		query,
		match.TeamAID,
		match.TeamBID,
		match.HostUserID,
		match.Location,
		match.MatchDate,
		match.OversPerInnings,
		match.Status,
	).Scan(
		&match.ID,
		&match.CreatedAt,
	)
}

func CreateMatchSnapshots(matchID string, teamAID string, teamBID string) (string, string, error) {
	_ = matchID
	return teamAID, teamBID, nil
}

func GetAllMatches(query dto.GetMatchesQuery) ([]dto.MatchResponse, error) {

	baseQuery := `
		SELECT
			m.id,
			m.host_user_id,
			m.host_user_id AS created_by,
			m.team_a_id,
			 ta_legacy.name AS team_a_name,
			m.team_b_id,
			tb_legacy.name AS team_b_name,
			COALESCE(score_summary.team_one_score, '—') AS team_one_score,
			COALESCE(score_summary.team_two_score, '—') AS team_two_score,
			m.team_a_id AS team_a_match_team_id,
			m.team_b_id AS team_b_match_team_id,
			m.location,
			m.status,
			m.match_date,
			m.overs_per_innings,
			m.toss_decision,
			m.toss_winner_team_id,
			m.winner_match_team_id
		FROM matches m
		LEFT JOIN teams ta_legacy ON ta_legacy.id = m.team_a_id
		LEFT JOIN teams tb_legacy ON tb_legacy.id = m.team_b_id
		LEFT JOIN LATERAL (
			SELECT
				MAX(CASE WHEN i.batting_team_id = m.team_a_id THEN CONCAT(COALESCE(i.total_runs, 0), '/', COALESCE(i.total_wickets, 0)) END) AS team_one_score,
				MAX(CASE WHEN i.batting_team_id = m.team_b_id THEN CONCAT(COALESCE(i.total_runs, 0), '/', COALESCE(i.total_wickets, 0)) END) AS team_two_score
			FROM innings i
			WHERE i.match_id = m.id
		) score_summary ON TRUE
		WHERE 1=1
	`

	var matches []dto.MatchResponse
	args := []interface{}{}

	if query.Status != "" {
		baseQuery += fmt.Sprintf(" AND m.status = $%d", len(args)+1)
		args = append(args, query.Status)
	}

	if query.TeamID != "" {
		baseQuery += fmt.Sprintf(" AND (m.team_a_id = $%d OR m.team_b_id = $%d)", len(args)+1, len(args)+2)
		args = append(args, query.TeamID, query.TeamID)
	}

	if query.HostUserID != "" {
		baseQuery += fmt.Sprintf(" AND m.host_user_id = $%d", len(args)+1)
		args = append(args, query.HostUserID)
	}

	if query.Search != "" {
		baseQuery += `
			AND (
				LOWER(COALESCE(ta_legacy.name, '')) LIKE $%d
				OR LOWER(COALESCE(tb_legacy.name, '')) LIKE $%d
				OR LOWER(m.location) LIKE $%d
			)
		`
		baseQuery = fmt.Sprintf(baseQuery, len(args)+1, len(args)+2, len(args)+3)
		search := "%" + strings.ToLower(query.Search) + "%"
		args = append(args, search, search, search)
	}

	baseQuery += " ORDER BY m.match_date DESC"

	err := database.DB.Select(&matches, baseQuery, args...)
	if err != nil {
		return nil, err
	}

	return matches, nil
}
func GetMatchByID(matchID string) (*models.Match, error) {

	var match models.Match

	query := `
	SELECT
		id,
		host_user_id,
		location,
		overs_per_innings,
		status,
		created_at,
		team_a_id,
		team_b_id
	FROM matches
	WHERE id = $1
	`

	err := database.DB.Get(&match, query, matchID)
	if err != nil {
		return nil, err
	}

	return &match, nil
}

func GetMatchDetailByID(matchID string) (*dto.MatchResponse, error) {
	var match dto.MatchResponse

	query := `
	SELECT
		m.id,
		m.host_user_id,
		m.host_user_id AS created_by,
		m.team_a_id,
		ta.name AS team_a_name,
		m.team_b_id,
		tb.name AS team_b_name,
		COALESCE(score_summary.team_one_score, '—') AS team_one_score,
		COALESCE(score_summary.team_two_score, '—') AS team_two_score,
		m.team_a_id AS team_a_match_team_id,
		m.team_b_id AS team_b_match_team_id,
		m.location,
		m.status,
		m.match_date,
		m.overs_per_innings,
		m.toss_decision,
		m.toss_winner_team_id,
		m.first_pick_team_id,
		m.winner_match_team_id
	FROM matches m
	LEFT JOIN teams ta ON ta.id = m.team_a_id
	LEFT JOIN teams tb ON tb.id = m.team_b_id
	LEFT JOIN LATERAL (
		SELECT
			MAX(CASE WHEN i.batting_team_id = m.team_a_id THEN CONCAT(COALESCE(i.total_runs, 0), '/', COALESCE(i.total_wickets, 0)) END) AS team_one_score,
			MAX(CASE WHEN i.batting_team_id = m.team_b_id THEN CONCAT(COALESCE(i.total_runs, 0), '/', COALESCE(i.total_wickets, 0)) END) AS team_two_score
		FROM innings i
		WHERE i.match_id = m.id
	) score_summary ON TRUE
	WHERE m.id = $1
	`

	if err := database.DB.Get(&match, query, matchID); err != nil {
		return nil, err
	}

	return &match, nil
}

func UpdateFirstPickTeam(matchID string, firstPickTeamID string) error {
	query := `
	UPDATE matches
	SET first_pick_team_id = $1
	WHERE id = $2
	`

	_, err := database.DB.Exec(query, firstPickTeamID, matchID)
	return err
}

func GetMatchInnings(matchID string) ([]dto.InningsResponse, error) {
	var innings []dto.InningsResponse

	query := `
	SELECT
		i.id,
		i.match_id,
		i.innings_no,
		i.batting_team_id,
		i.bowling_team_id,
		COALESCE(i.total_runs, 0) AS total_runs,
		COALESCE(i.total_wickets, 0) AS total_wickets,
		COALESCE(s.legal_balls, legal_totals.legal_balls, 0) AS legal_balls,
		COALESCE(s.current_over, FLOOR(COALESCE(s.legal_balls, legal_totals.legal_balls, 0) / 6.0)::INT) AS current_over,
		COALESCE(s.current_ball, MOD(COALESCE(s.legal_balls, legal_totals.legal_balls, 0), 6)) AS current_ball
	FROM innings i
	LEFT JOIN innings_state s ON s.innings_id = i.id
	LEFT JOIN LATERAL (
		SELECT COALESCE(SUM(CASE WHEN COALESCE(be.ball_type, 'normal') NOT IN ('wide','no_ball','dead_ball','retired_hurt') THEN 1 ELSE 0 END), 0)::INT AS legal_balls
		FROM ball_events be
		WHERE be.innings_id = i.id
		  AND be.is_deleted = FALSE
	) legal_totals ON TRUE
	WHERE i.match_id = $1
	ORDER BY i.innings_no ASC`

	if err := database.DB.Select(&innings, query, matchID); err != nil {
		return nil, err
	}

	return innings, nil
}

func GetMatchScorecard(matchID string) (*dto.MatchScorecardResponse, error) {
	scorecard := &dto.MatchScorecardResponse{}

	innings, err := GetMatchInnings(matchID)
	if err != nil {
		return nil, err
	}
	scorecard.Innings = innings

	battingQuery := `
	SELECT
		pms.team_player_id,
		tp.team_id,
		tp.player_id AS user_id,
		COALESCE(u.name, '') AS player_name,
		COALESCE(pms.runs_scored, 0) AS runs_scored,
		COALESCE(pms.balls_faced, 0) AS balls_faced,
		COALESCE(pms.fours, 0) AS fours,
		COALESCE(pms.sixes, 0) AS sixes,
		COALESCE(pms.is_out, FALSE) AS is_out,
		COALESCE(pms.runs_conceded, 0) AS runs_conceded,
		COALESCE(pms.wickets_taken, 0) AS wickets_taken,
		(
			FLOOR(COALESCE(pms.legal_balls_bowled, 0) / 6.0)
			+ MOD(COALESCE(pms.legal_balls_bowled, 0), 6)::NUMERIC / 10.0
		)::FLOAT AS overs_bowled,
		COALESCE(pms.fantasy_points, 0) AS fantasy_points
	FROM player_match_stats pms
	JOIN team_players tp ON tp.id = pms.team_player_id
	LEFT JOIN users u ON u.id = tp.player_id
	WHERE pms.match_id = $1
	ORDER BY pms.runs_scored DESC, pms.balls_faced DESC
	`
	if err := database.DB.Select(&scorecard.Batting, battingQuery, matchID); err != nil {
		return nil, err
	}

	bowlingQuery := `
	SELECT
		pms.team_player_id,
		tp.team_id,
		tp.player_id AS user_id,
		COALESCE(u.name, '') AS player_name,
		COALESCE(pms.runs_scored, 0) AS runs_scored,
		COALESCE(pms.balls_faced, 0) AS balls_faced,
		COALESCE(pms.fours, 0) AS fours,
		COALESCE(pms.sixes, 0) AS sixes,
		COALESCE(pms.is_out, FALSE) AS is_out,
		COALESCE(pms.runs_conceded, 0) AS runs_conceded,
		COALESCE(pms.wickets_taken, 0) AS wickets_taken,
		(
			FLOOR(COALESCE(pms.legal_balls_bowled, 0) / 6.0)
			+ MOD(COALESCE(pms.legal_balls_bowled, 0), 6)::NUMERIC / 10.0
		)::FLOAT AS overs_bowled,
		COALESCE(pms.fantasy_points, 0) AS fantasy_points
	FROM player_match_stats pms
	JOIN team_players tp ON tp.id = pms.team_player_id
	LEFT JOIN users u ON u.id = tp.player_id
	WHERE pms.match_id = $1
	ORDER BY pms.wickets_taken DESC, COALESCE(pms.legal_balls_bowled, 0) DESC
	`
	if err := database.DB.Select(&scorecard.Bowling, bowlingQuery, matchID); err != nil {
		return nil, err
	}

	recentBallsQuery := `
	SELECT
		id,
		innings_id,
		ball_no,
		COALESCE(delivery_no, delivery_number, 1) AS delivery_no,
		COALESCE(ball_type, 'normal') AS ball_type,
		COALESCE(total_runs, 0) AS total_runs,
-- 		COALESCE(is_wicket, FALSE) AS is_wicket,
		striker_id,
		non_striker_id,
		bowler_id,
		dismissal_type
	FROM ball_events
	WHERE match_id = $1
	  AND is_deleted = FALSE
	ORDER BY created_at DESC
	LIMIT 12
	`
	if err := database.DB.Select(&scorecard.RecentBalls, recentBallsQuery, matchID); err != nil {
		return nil, err
	}

	var current struct {
		StrikerID    *string `db:"striker_id"`
		NonStrikerID *string `db:"non_striker_id"`
		BowlerID     *string `db:"bowler_id"`
	}
	err = database.DB.Get(&current, `
		SELECT s.striker_id, s.non_striker_id, s.bowler_id
		FROM innings i
		LEFT JOIN innings_state s ON s.innings_id = i.id
		WHERE i.match_id = $1
		ORDER BY i.innings_no DESC
		LIMIT 1
	`, matchID)
	if err == nil {
		scorecard.CurrentStrikerID = current.StrikerID
		scorecard.CurrentNonStrikerID = current.NonStrikerID
		scorecard.CurrentBowlerID = current.BowlerID
	}

	return scorecard, nil
}

func ResolveMatchTeamID(matchID string, teamOrMatchTeamID string) (string, error) {
	var resolvedID string

	query := `
	SELECT id
	FROM teams
	WHERE (id::text = $2)
	  AND id IN (
		SELECT team_a_id FROM matches WHERE id = $1
		UNION
		SELECT team_b_id FROM matches WHERE id = $1
	  )
	LIMIT 1
	`

	if err := database.DB.Get(&resolvedID, query, matchID, teamOrMatchTeamID); err != nil {
		return "", err
	}

	return resolvedID, nil
}

func StartMatch(matchID string) ([]dto.InningsResponse, error) {
	tx, err := database.DB.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
		UPDATE matches
		SET status = 'live', started_at = COALESCE(started_at, NOW())
		WHERE id = $1
	`, matchID)
	if err != nil {
		return nil, err
	}

	var inningsCount int
	if err := tx.Get(&inningsCount, `SELECT COUNT(1) FROM innings WHERE match_id = $1`, matchID); err != nil {
		return nil, err
	}

	if inningsCount == 0 {
		var teamA, teamB string
		if err := tx.QueryRowx(`SELECT team_a_id, team_b_id FROM matches WHERE id = $1`, matchID).Scan(&teamA, &teamB); err != nil {
			return nil, err
		}

		if _, err := tx.Exec(`
			INSERT INTO innings (
				match_id,
				batting_team_id,
				bowling_team_id,
				innings_no,
				started_at
			) VALUES ($1, $2, $3, 1, NOW())
		`, matchID, teamA, teamB); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return GetMatchInnings(matchID)
}

func GetMatchSquad(matchID string) ([]dto.MatchSquadPlayer, error) {
	var players []dto.MatchSquadPlayer

	query := `
	SELECT
		tp.id AS team_player_id,
		tp.team_id,
		tp.player_id AS user_id,
		COALESCE(u.name, '') AS player_name,
		u.phone_number,
		COALESCE(tp.is_playing_xi, FALSE) AS is_playing_xi,
		COALESCE(tp.is_captain, FALSE) AS is_captain,
		COALESCE(tp.is_umpire, FALSE) AS is_umpire,
		tp.batting_order
	FROM team_players tp
	LEFT JOIN users u ON u.id = tp.player_id
	WHERE tp.team_id IN (
		SELECT team_a_id FROM matches WHERE id = $1
		UNION
		SELECT team_b_id FROM matches WHERE id = $1
	)
	  AND tp.deleted_at IS NULL
	ORDER BY tp.team_id ASC, tp.batting_order ASC NULLS LAST, tp.created_at ASC
	`

	if err := database.DB.Select(&players, query, matchID); err != nil {
		return nil, err
	}

	return players, nil
}

func UpdateMatchLineup(matchID string, players []dto.UpdateLineupPlayer) error {
	tx, err := database.DB.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, p := range players {
		_, err := tx.Exec(`
			UPDATE team_players tp
			SET
				is_playing_xi = $1,
				is_captain = $2,
				is_umpire = $3,
				batting_order = $4,
				updated_at = NOW()
			WHERE tp.id = $5
			  AND tp.team_id IN (
				SELECT team_a_id FROM matches WHERE id = $6
				UNION
				SELECT team_b_id FROM matches WHERE id = $6
			  )
		`, p.IsPlayingXI, p.IsCaptain, p.IsUmpire, p.BattingOrder, p.MatchTeamPlayerID, matchID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func GetMatchInningsIDs(matchID string) ([]string, error) {
	var ids []string
	err := database.DB.Select(&ids, `
		SELECT id
		FROM innings
		WHERE match_id = $1
		ORDER BY innings_no ASC`, matchID)
	if err != nil {
		return nil, err
	}
	return ids, nil
}

func FinalizeMatch(matchID, winnerMatchTeamID string) error {
	tx, err := database.DB.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var pomUserID *string
	var worstUserID *string

	err = tx.QueryRowx(`
		SELECT tp.player_id
		FROM player_match_stats pms
		JOIN team_players tp ON tp.id = pms.team_player_id
		WHERE pms.match_id = $1
		  AND tp.player_id IS NOT NULL
		ORDER BY pms.fantasy_points DESC, pms.updated_at DESC
		LIMIT 1
	`, matchID).Scan(&pomUserID)
	if err != nil {
		pomUserID = nil
	}

	err = tx.QueryRowx(`
		SELECT tp.player_id
		FROM player_match_stats pms
		JOIN team_players tp ON tp.id = pms.team_player_id
		WHERE pms.match_id = $1
		  AND tp.player_id IS NOT NULL
		ORDER BY pms.fantasy_points ASC, pms.updated_at DESC
		LIMIT 1
	`, matchID).Scan(&worstUserID)
	if err != nil {
		worstUserID = nil
	}

	_, err = tx.Exec(`
		UPDATE matches
		SET
			winner_match_team_id = $2,
			status = 'completed',
			completed_at = NOW(),
			player_of_match_user_id = $3,
			worst_player_user_id = $4
		WHERE id = $1
	`, matchID, winnerMatchTeamID, pomUserID, worstUserID)
	if err != nil {
		return err
	}

	return tx.Commit()
}
