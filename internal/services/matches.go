package services

import (
	"CricketDuniya-Backend/internal/dto"
	"CricketDuniya-Backend/internal/models"
	"CricketDuniya-Backend/internal/repositories"
	"time"
)

func CreateMatch(req dto.CreateMatchRequest, hostUserID string) (*models.Match, error) {

	var parsedDate *time.Time

	if req.MatchDate != "" {

		t, err := time.Parse(time.RFC3339, req.MatchDate)
		if err != nil {
			return nil, err
		}

		parsedDate = &t
	}

	match := &models.Match{
		TeamAID:         &req.TeamAID,
		TeamBID:         &req.TeamBID,
		HostUserID:      hostUserID,
		Location:        &req.Location,
		MatchDate:       parsedDate,
		OversPerInnings: req.OversPerInnings,
		Status:          "scheduled",
	}

	err := repositories.CreateMatch(match)
	if err != nil {
		return nil, err
	}

	return match, nil
}

func GetAllMatches(query dto.GetMatchesQuery) ([]dto.MatchResponse, error) {

	return repositories.GetAllMatches(query)
}

func GetMatchByID(id string) (*models.Match, error) {
	return repositories.GetMatchByID(id)
}
