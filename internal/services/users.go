package services

import (
	"CricketDuniya-Backend/internal/dto"
	"CricketDuniya-Backend/internal/repositories"
)

func GetAllUsers() ([]dto.UserProfileUser, error) {
	return repositories.GetAllUsers()
}
