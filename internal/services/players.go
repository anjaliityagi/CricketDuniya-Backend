package services

import (
	"CricketDuniya-Backend/internal/dto"
	"CricketDuniya-Backend/internal/models"
	"CricketDuniya-Backend/internal/repositories"
	"errors"

	"github.com/google/uuid"
)

var ErrPlayerNotFound = errors.New("player not found")

func CreatePlayer(req dto.CreatePlayerRequest, teamID uuid.UUID) (*models.MatchPlayer, error) {

	player := &models.MatchPlayer{

		PlayerID:  req.PLayerID,
		IsHost:    false,
		IsCaptain: false,
	}

	err := repositories.CreatePlayer(player, teamID)
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
