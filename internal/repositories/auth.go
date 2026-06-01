package repositories

import (
	"CricketDuniya-Backend/internal/database"
	"CricketDuniya-Backend/internal/dto"
	"CricketDuniya-Backend/internal/models"
	"database/sql"
	"errors"
	"fmt"
)

func CreateUser(user *models.User) error {

	query := `
	INSERT INTO users (
		name,
		phone_number,
		password_hash
	)
	VALUES ($1, $2, $3)
	RETURNING id, created_at, updated_at
	`

	return database.DB.QueryRow(
		query,
		user.Name,
		user.PhoneNumber,
		user.PasswordHash,
	).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
}

func CreateGuestUser(user *models.User) error {

	query := `
	INSERT INTO users (
		name,
		phone_number,
		password_hash,
		is_phone_verified
	)
	VALUES ($1, $2, NULL, FALSE)
	RETURNING id, created_at, updated_at
	`

	return database.DB.QueryRow(
		query,
		user.Name,
		user.PhoneNumber,
	).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
}
func GetUserByPhoneNumber(phoneNumber string) (*models.User, error) {

	var user models.User

	query := `
	SELECT
		id,
		name,
		phone_number,
		password_hash,
		COALESCE(is_phone_verified, FALSE) AS is_phone_verified,
		created_at,
		updated_at
	FROM users
	WHERE phone_number = $1
	`

	err := database.DB.Get(&user, query, phoneNumber)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		fmt.Println(err)
		return nil, err
	}

	return &user, nil
}

func CreateSession(userID string) (*dto.UserSession, error) {

	query := `
	INSERT INTO user_session (user_id)
	VALUES ($1)
	RETURNING id, user_id, created_at, archived_at
	`

	var session dto.UserSession

	err := database.DB.QueryRowx(query, userID).StructScan(&session)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func LogoutSession(sessionID string) error {

	query := `
	UPDATE user_session
	SET archived_at = NOW()
	WHERE id = $1
	`

	_, err := database.DB.Exec(query, sessionID)
	return err
}

func IsSessionActive(sessionID string) (bool, error) {

	query := `
	SELECT id
	FROM user_session
	WHERE id = $1
	AND archived_at IS NULL
	`

	var id string

	err := database.DB.Get(&id, query, sessionID)
	if err != nil {
		return false, nil
	}

	return true, nil
}
