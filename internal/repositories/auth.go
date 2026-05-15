package repositories

import (
	"CricketDuniya-Backend/internal/database"
	"CricketDuniya-Backend/internal/models"
	"fmt"
)

func CreateUser(user *models.User) error {

	query := `
	INSERT INTO users (
		name,
		email,
		password_hash
	)
	VALUES ($1, $2, $3)
	RETURNING id, created_at, updated_at
	`

	return database.DB.QueryRow(
		query,
		user.Name,
		user.Email,
		user.PasswordHash,
	).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
}

func GetUserByEmail(email string) (*models.User, error) {

	var user models.User

	query := `SELECT
		id,
		name,
		email,
		password_hash,
		profile_image,
		created_at,
		updated_at
	FROM users
	WHERE email = $1`

	err := database.DB.Get(&user, query, email)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &user, nil
}
