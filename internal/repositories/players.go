package repositories

import (
	"CricketDuniya-Backend/internal/database"
	"CricketDuniya-Backend/internal/models"
	"fmt"

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
	VALUES ($1, $2) ON CONFLICT ON CONSTRAINT unique_team_user_match
DO NOTHING
	`

	_, err := database.DB.Exec(
		query,
		teamPlayer.TeamID,
		teamPlayer.PlayerID,
	)
	fmt.Println(err)
	return err
}

func AssignCaptain(
	teamID string,
	playerID string,
) error {

	resetQuery := `
	UPDATE team_players
	SET is_captain = false
	WHERE team_id = $1
	`

	_, err := database.DB.Exec(
		resetQuery,
		teamID,
	)

	if err != nil {
		return err
	}

	assignQuery := `
	UPDATE team_players
	SET is_captain = true
	WHERE team_id = $1
	AND player_id = $2
	`

	_, err = database.DB.Exec(
		assignQuery,
		teamID,
		playerID,
	)

	return err
}

func AssignUmpire(
	teamID string,
	playerID string,
) error {

	resetQuery := `
	UPDATE team_players
	SET is_umpire = false
	WHERE team_id = $1
	`

	_, err := database.DB.Exec(
		resetQuery,
		teamID,
	)

	if err != nil {
		return err
	}

	assignQuery := `
	UPDATE team_players
	SET is_umpire = true
	WHERE team_id = $1
	AND player_id = $2
	`

	_, err = database.DB.Exec(
		assignQuery,
		teamID,
		playerID,
	)

	return err
}
