package repositories

import (
	"CricketDuniya-Backend/internal/database"
	"CricketDuniya-Backend/internal/dto"
	"CricketDuniya-Backend/internal/models"
	"strings"
)

func CreateMatch(match *models.Match) error {

	query := `
	INSERT INTO matches (
		host_user_id,
		venue,
		match_date,
		overs_per_side,
		status
	)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, created_at
	`

	return database.DB.QueryRowx(
		query,
		match.HostUserID,
		match.Venue,
		match.MatchDate,
		match.OversPerSide,
		match.Status,
	).Scan(
		&match.ID,
		&match.CreatedAt,
	)
}
func GetAllMatches(query dto.GetMatchesQuery) ([]dto.MatchResponse, error) {

	baseQuery := `
		SELECT
			m.id,
			m.team_a_id,
			ta.name AS team_a_name,
			m.team_b_id,
			tb.name AS team_b_name,
			m.venue,
			m.status,
			m.match_date,
			m.overs_per_side,
			m.toss_decision,
			m.winner_team_id
		FROM matches m
		left join  teams ta ON ta.id = m.team_a_id
		left join teams tb ON tb.id = m.team_b_id
		WHERE 1=1
	`

	var matches []dto.MatchResponse

	if query.Status != "" {
		baseQuery += " AND m.status = ?"
	}

	if query.TeamID != "" {
		baseQuery += " AND (m.team_a_id = ? OR m.team_b_id = ?)"
	}

	if query.HostUserID != "" {
		baseQuery += " AND m.host_user_id = ?"
	}

	if query.Search != "" {
		baseQuery += `
			AND (
				LOWER(ta.name) LIKE ?
				OR LOWER(tb.name) LIKE ?
				OR LOWER(m.venue) LIKE ?
			)
		`
	}

	baseQuery += " ORDER BY m.match_date DESC"

	// build args simply
	args := []interface{}{}

	if query.Status != "" {
		args = append(args, query.Status)
	}

	if query.TeamID != "" {
		args = append(args, query.TeamID, query.TeamID)
	}

	if query.HostUserID != "" {
		args = append(args, query.HostUserID)
	}

	if query.Search != "" {
		search := "%" + strings.ToLower(query.Search) + "%"
		args = append(args, search, search, search)
	}

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
		venue,
		overs_per_side,
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
