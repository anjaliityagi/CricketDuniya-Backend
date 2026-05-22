package services

import (
	"CricketDuniya-Backend/internal/dto"
	"CricketDuniya-Backend/internal/repositories"
)

func GetAllUsers(search string) ([]dto.UserProfileUser, error) {
	return repositories.GetAllUsers(search)
}
