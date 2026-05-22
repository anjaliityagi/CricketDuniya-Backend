package repositories

import (
	"CricketDuniya-Backend/internal/database"
	"CricketDuniya-Backend/internal/dto"
	"strings"
)

func GetAllUsers(search string) ([]dto.UserProfileUser, error) {
	var users []dto.UserProfileUser

	query := `SELECT
		id,
		name,
		phone_number,
		created_at,
		updated_at
	FROM users`

	args := []interface{}{}
	if strings.TrimSpace(search) != "" {
		query += `
		WHERE phone_number ILIKE $1
		   OR name ILIKE $1`
		args = append(args, "%"+strings.TrimSpace(search)+"%")
	}
	query += `
	ORDER BY created_at DESC
	LIMIT 50`

	err := database.DB.Select(&users, query, args...)
	if err != nil {
		return nil, err
	}

	return users, nil
}
