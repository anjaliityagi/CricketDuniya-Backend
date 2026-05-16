package models

import "time"

type User struct {
	ID           string    `db:"id" json:"id"`
	Name         string    `db:"name" json:"name"`
	PhoneNumber  string    `db:"phone_number" json:"phone_number"`
	PasswordHash string    `db:"password_hash" json:"-"`
	ProfileImage *string   `db:"profile_image" json:"profile_image"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}
