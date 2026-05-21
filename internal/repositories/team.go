package repositories

import (
	"CricketDuniya-Backend/internal/database"
	"CricketDuniya-Backend/internal/models"
)

func CreateTeam(team *models.Team) error {

	query := `
	INSERT INTO teams (name, created_by)
	VALUES ($1, $2)
	RETURNING id, created_at
	`

	return database.DB.QueryRow(
		query,
		team.Name,
		team.CreatedBy,
	).Scan(&team.ID, &team.CreatedAt)
}

func GetTeamByID(id string) (*models.Team, error) {

	var team models.Team

	query := `
	SELECT id, name, captain_id, created_by, created_at
	FROM teams
	WHERE id = $1
	`

	err := database.DB.Get(&team, query, id)
	if err != nil {
		return nil, err
	}

	return &team, nil
}

func GetAllTeams() ([]models.Team, error) {

	var teams []models.Team

	query := `
	SELECT id, name, captain_id, created_by, created_at
	FROM teams
	ORDER BY created_at DESC
	`

	err := database.DB.Select(&teams, query)
	if err != nil {
		return nil, err
	}

	return teams, nil
}
