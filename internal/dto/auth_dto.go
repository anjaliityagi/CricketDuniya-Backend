package dto

import "time"

type SignupRequest struct {
	Name        string `json:"name" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required,min=10,max=15"`
	Password    string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required,min=10,max=15"`
	Password    string `json:"password" binding:"required"`
}

type UserSession struct {
	ID         string     `db:"id"`
	UserID     string     `db:"user_id"`
	CreatedAt  time.Time  `db:"created_at"`
	ArchivedAt *time.Time `db:"archived_at"`
}
