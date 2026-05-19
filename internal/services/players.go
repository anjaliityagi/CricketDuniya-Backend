package services

import (
	"CricketDuniya-Backend/internal/dto"
	"CricketDuniya-Backend/internal/models"
	"CricketDuniya-Backend/internal/repositories"
	"errors"
)

var ErrPlayerNotFound = errors.New("player not found")

func CreatePlayer(req dto.CreatePlayerRequest, userID string) (*models.MatchPlayer, error) {

	player := &models.MatchPlayer{
		MatchID:        req.MatchID,
		UserID:         &userID,
		PlayerName:     req.PlayerName,
		Phone:          req.Phone,
		TeamID:         req.TeamID,
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

func DeletePlayer(playerID string, userID string) error {
	deleted, err := repositories.DeletePlayer(playerID, userID)
	if err != nil {
		return err
	}

	if !deleted {
		return ErrPlayerNotFound
	}

	return nil
}
