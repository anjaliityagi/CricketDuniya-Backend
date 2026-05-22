package models

import "time"

type User struct {
	ID              string    `db:"id" json:"id"`
	Name            string    `db:"name" json:"name"`
	PhoneNumber     string    `db:"phone_number" json:"phone_number"`
	PasswordHash    *string   `db:"password_hash" json:"-"`
	IsPhoneVerified bool      `db:"is_phone_verified" json:"is_phone_verified"`
	BattingStyle    *string   `db:"batting_style" json:"batting_style"`
	BowlingStyle    *string   `db:"bowling_style" json:"bowling_style"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time `db:"updated_at" json:"updated_at"`
}
