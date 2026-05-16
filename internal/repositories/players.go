package repositories

import (
	"CricketDuniya-Backend/internal/database"
	"CricketDuniya-Backend/internal/models"
)

func CreatePlayer(player *models.MatchPlayer) error {

	query := `
	INSERT INTO match_players (
		match_id,
		user_id,
		player_name,
		phone,
		team_side,
		is_host,
		is_captain,
		is_wicketkeeper
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	RETURNING id, created_at
	`

	return database.DB.QueryRowx(
		query,
		player.MatchID,
		player.UserID,
		player.PlayerName,
		player.Phone,
		player.TeamSide,
		player.IsHost,
		player.IsCaptain,
		player.IsWicketkeeper,
	).Scan(
		&player.ID,
		&player.CreatedAt,
	)
}
