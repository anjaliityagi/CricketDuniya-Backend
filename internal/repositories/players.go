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
	ON CONFLICT (team_id, player_id)
	DO UPDATE SET
		is_captain = EXCLUDED.is_captain,
		deleted_at = NULL,
		updated_at = NOW()
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
	UPDATE team_players
	SET deleted_at = NOW(), updated_at = NOW()
	WHERE id = $1
	AND deleted_at IS NULL
	`

	_ = userID
	result, err := database.DB.Exec(query, playerID)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected > 0, nil
}

func AddPlayerToTeam(
	teamPlayer *models.TeamPlayer,
) error {

	query := `
	INSERT INTO team_players (
		team_id,
		player_id
	)
	VALUES ($1, $2)
	`

	_, err := database.DB.Exec(
		query,
		teamPlayer.TeamID,
		teamPlayer.PlayerID,
	)

	return err
}
