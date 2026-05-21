package repositories

import (
	"CricketDuniya-Backend/internal/database"
	"CricketDuniya-Backend/internal/dto"
)

func GetAllUsers() ([]dto.UserProfileUser, error) {
	var users []dto.UserProfileUser

	query := `SELECT
		id,
		name,
		phone_number,
		created_at,
		updated_at
	FROM users
	ORDER BY created_at DESC`

	err := database.DB.Select(&users, query)
	if err != nil {
		return nil, err
	}

	return users, nil
}
