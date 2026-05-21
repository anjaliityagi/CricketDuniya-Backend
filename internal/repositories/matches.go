package repositories

import (
	"CricketDuniya-Backend/internal/database"
	"CricketDuniya-Backend/internal/dto"
	"CricketDuniya-Backend/internal/models"
	"fmt"
	"strings"
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

func CreateMatchSnapshots(matchID string, teamAID string, teamBID string) (string, string, error) {
	tx, err := database.DB.Beginx()
	if err != nil {
		return "", "", err
	}
	defer tx.Rollback()

	createSnapshot := func(sourceTeamID string) (string, error) {
		var matchTeamID string
		insertSnapshot := `
			INSERT INTO match_teams (match_id, team_id)
			SELECT $1, t.id
			FROM teams t
			WHERE t.id = $2
			RETURNING id
		`
		if err := tx.QueryRow(insertSnapshot, matchID, sourceTeamID).Scan(&matchTeamID); err != nil {
			return "", err
		}

		copyPlayers := `
			INSERT INTO match_team_players (
				match_team_id,
				user_id,
				player_name,
				phone_number,
				is_playing_xi,
				is_substitute,
				is_captain,
				is_wicket_keeper,
				batting_order
			)
			SELECT
				$1,
				tp.player_id,
				u.name,
				u.phone_number,
				TRUE,
				COALESCE(tp.is_substitute, FALSE),
				COALESCE(tp.is_captain, FALSE),
				COALESCE(tp.is_wicket_keeper, FALSE),
				tp.batting_order
			FROM team_players tp
			JOIN users u ON u.id = tp.player_id
			WHERE tp.team_id = $2
			  AND tp.removed_at IS NULL
		`
		if _, err := tx.Exec(copyPlayers, matchTeamID, sourceTeamID); err != nil {
			return "", err
		}

		return matchTeamID, nil
	}

	matchTeamAID, err := createSnapshot(teamAID)
	if err != nil {
		return "", "", err
	}

	matchTeamBID, err := createSnapshot(teamBID)
	if err != nil {
		return "", "", err
	}

	if err := tx.Commit(); err != nil {
		return "", "", err
	}

	return matchTeamAID, matchTeamBID, nil
}

func GetAllMatches(query dto.GetMatchesQuery) ([]dto.MatchResponse, error) {

	baseQuery := `
		SELECT
			m.id,
			COALESCE(mt_a.id, m.team_a_id) AS team_a_id,
			 ta_legacy.name AS team_a_name,
			COALESCE(mt_b.id, m.team_b_id) AS team_b_id,
			tb_legacy.name AS team_b_name,
			m.location,
			m.status,
			m.match_date,
			m.overs_per_innings,
			m.toss_decision,
			m.winner_match_team_id
		FROM matches m
		LEFT JOIN LATERAL (
			SELECT id
			FROM match_teams
			WHERE match_id = m.id
			  AND deleted_at IS NULL
			ORDER BY created_at ASC
			LIMIT 1
		) mt_a ON TRUE
		LEFT JOIN LATERAL (
			SELECT id
			FROM match_teams
			WHERE match_id = m.id
			  AND deleted_at IS NULL
			ORDER BY created_at ASC
			OFFSET 1
			LIMIT 1
		) mt_b ON TRUE
		LEFT JOIN teams ta_legacy ON ta_legacy.id = m.team_a_id
		LEFT JOIN teams tb_legacy ON tb_legacy.id = m.team_b_id
		WHERE 1=1
	`

	var matches []dto.MatchResponse
	args := []interface{}{}

	if query.Status != "" {
		baseQuery += fmt.Sprintf(" AND m.status = $%d", len(args)+1)
		args = append(args, query.Status)
	}

	if query.TeamID != "" {
		baseQuery += fmt.Sprintf(" AND (m.team_a_id = $%d OR m.team_b_id = $%d OR mt_a.id = $%d OR mt_b.id = $%d)", len(args)+1, len(args)+2, len(args)+3, len(args)+4)
		args = append(args, query.TeamID, query.TeamID, query.TeamID, query.TeamID)
	}

	if query.HostUserID != "" {
		baseQuery += fmt.Sprintf(" AND m.host_user_id = $%d", len(args)+1)
		args = append(args, query.HostUserID)
	}

	if query.Search != "" {
		baseQuery += `
			AND (
				LOWER(COALESCE(mt_a.display_name, ta_legacy.name, '')) LIKE $%d
				OR LOWER(COALESCE(mt_b.display_name, tb_legacy.name, '')) LIKE $%d
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

func GetMatchInningsIDs(matchID string) ([]string, error) {
	var ids []string
	err := database.DB.Select(&ids, `
		SELECT id
		FROM innings
		WHERE match_id = $1
		ORDER BY created_at ASC
	`, matchID)
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
		SELECT mtp.user_id
		FROM player_match_stats pms
		JOIN match_team_players mtp ON mtp.id = pms.match_team_player_id
		WHERE pms.match_id = $1
		  AND mtp.user_id IS NOT NULL
		ORDER BY pms.fantasy_points DESC, pms.updated_at DESC
		LIMIT 1
	`, matchID).Scan(&pomUserID)
	if err != nil {
		pomUserID = nil
	}

	err = tx.QueryRowx(`
		SELECT mtp.user_id
		FROM player_match_stats pms
		JOIN match_team_players mtp ON mtp.id = pms.match_team_player_id
		WHERE pms.match_id = $1
		  AND mtp.user_id IS NOT NULL
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
