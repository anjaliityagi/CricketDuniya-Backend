package services

import (
	"CricketDuniya-Backend/internal/dto"
	"CricketDuniya-Backend/internal/models"
	"CricketDuniya-Backend/internal/repositories"
)

func CreatePlayer(req dto.CreatePlayerRequest, userID string) (*models.MatchPlayer, error) {

	player := &models.MatchPlayer{
		MatchID:        req.MatchID,
		UserID:         &userID,
		PlayerName:     req.PlayerName,
		Phone:          req.Phone,
		TeamSide:       req.TeamSide,
		IsHost:         false,
		IsCaptain:      false,
		IsWicketkeeper: false,
	}

	err := repositories.CreatePlayer(player)
	if err != nil {
		return nil, err
	}

	return player, nil
}
