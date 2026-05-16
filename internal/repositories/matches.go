package repositories

import (
	"CricketDuniya-Backend/internal/database"
	"CricketDuniya-Backend/internal/models"
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
