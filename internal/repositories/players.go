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
		team_id,
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
		player.TeamID,
		player.IsHost,
		player.IsCaptain,
		player.IsWicketkeeper,
	).Scan(
		&player.ID,
		&player.CreatedAt,
	)
}

func DeletePlayer(playerID string, userID string) (bool, error) {
	query := `
	UPDATE match_players
	SET removed_at = NOW()
	WHERE id = $1
	AND user_id = $2
	AND removed_at IS NULL
	`

	result, err := database.DB.Exec(query, playerID, userID)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected > 0, nil
}
