package repositories

import (
	"CricketDuniya-Backend/internal/database"
	"CricketDuniya-Backend/internal/models"

	"github.com/google/uuid"
)

func CreatePlayer(player *models.MatchPlayer, teamID uuid.UUID) error {

	query := `
	INSERT INTO team_players (
		
		player_id,
	
		
		team_id,
		
		is_captain
	
	)
	VALUES ($1, $2, $3)
	RETURNING id,team_id, created_at
	`

	return database.DB.QueryRowx(
		query,

		player.PlayerID,
		teamID,

		player.IsCaptain,
	).Scan(
		&player.ID,
		&player.TeamID,
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
