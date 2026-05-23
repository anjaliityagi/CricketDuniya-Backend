package services

import (
	"CricketDuniya-Backend/internal/dto"
	"CricketDuniya-Backend/internal/models"
	"CricketDuniya-Backend/internal/repositories"
	"fmt"
)

func CreateTeam(req dto.CreateTeamRequest, userID string) (*models.Team, error) {

	team := &models.Team{
		Name:      &req.Name,
		CreatedBy: &userID,
	}
	err := repositories.CreateTeam(team)
	fmt.Println(team.ID)
	if err != nil {
		return nil, err
	}
	return team, nil
}

func GetTeamById(id string) (*models.Team, error) {
	return repositories.GetTeamByID(id)
}

func GetTeams() ([]models.Team, error) {
	return repositories.GetAllTeams()
}
