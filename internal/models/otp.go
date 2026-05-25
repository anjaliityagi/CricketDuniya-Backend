package models

import "time"

type OTP struct {
	ID        string    `db:"id"`
	Phone     string    `db:"phone"`
	OTPCode   string    `db:"otp_code"`
	ExpiresAt time.Time `db:"expires_at"`
	Used      bool      `db:"used"`
}
