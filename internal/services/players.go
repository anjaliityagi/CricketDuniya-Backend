package services

import (
	"CricketDuniya-Backend/internal/dto"
	"CricketDuniya-Backend/internal/models"
	"CricketDuniya-Backend/internal/repositories"
	"errors"
	"strings"

	"github.com/google/uuid"
)

var ErrPlayerNotFound = errors.New("player not found")

func CreatePlayer(req dto.CreatePlayerRequest, teamID uuid.UUID) (*models.MatchPlayer, error) {
	var resolvedUserID *string

	if req.PlayerID != nil && strings.TrimSpace(*req.PlayerID) != "" {
		id := strings.TrimSpace(*req.PlayerID)
		resolvedUserID = &id
	} else {
		if req.Name == nil || strings.TrimSpace(*req.Name) == "" {
			return nil, errors.New("name is required when player_id is not provided")
		}
		if req.PhoneNumber == nil || strings.TrimSpace(*req.PhoneNumber) == "" {
			return nil, errors.New("phone_number is required when player_id is not provided")
		}

		user, err := repositories.GetOrCreateLiteUserByPhone(strings.TrimSpace(*req.Name), strings.TrimSpace(*req.PhoneNumber))
		if err != nil {
			return nil, err
		}
		resolvedUserID = &user.ID
	}

	player := &models.MatchPlayer{
		PlayerID:  req.PlayerID,
		IsHost:    false,
		IsCaptain: false,
	}
	player.PlayerID = resolvedUserID

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

func GetPlayersByTeam(teamID string) ([]dto.TeamPlayerResponse, error) {
	return repositories.GetPlayersByTeamID(teamID)
}
