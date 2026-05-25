package repositories

import (
	"CricketDuniya-Backend/internal/database"
	"time"
)

func SaveOTP(phone string, otp string) error {

	query := `
	INSERT INTO otp_codes (
		phone,
		otp_code,
		expires_at,
		used
	)
	VALUES ($1, $2, $3, false)
	`

	_, err := database.DB.Exec(
		query,
		phone,
		otp,
		time.Now().Add(5*time.Minute),
	)

	return err
}

func VerifyOTP(phone string, otp string) (bool, error) {

	query := `
	SELECT COUNT(1)
	FROM otp_codes
	WHERE phone = $1
	  AND otp_code = $2
	  AND used = false
	  AND expires_at > NOW()
	`

	var count int

	err := database.DB.Get(&count, query, phone, otp)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func MarkOTPUsed(phone string, otp string) error {

	query := `
	UPDATE otp_codes
	SET used = true
	WHERE phone = $1
	  AND otp_code = $2
	`

	_, err := database.DB.Exec(query, phone, otp)

	return err
}

func UpdatePassword(phone string, hashedPassword string) error {

	query := `
	UPDATE users
	SET password_hash = $1
	WHERE phone_number = $2
	`

	_, err := database.DB.Exec(query, hashedPassword, phone)

	return err
}
