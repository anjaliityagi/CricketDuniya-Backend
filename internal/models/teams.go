package models

type Team struct {
	ID        *string `json:"id" db:"id"`
	Name      *string `json:"name" db:"name"`
	LogoURL   *string `json:"logo_url" db:"logo_url"`
	CaptainID *string `json:"captain_id" db:"captain_id"`
	CreatedBy *string `json:"created_by" db:"created_by"`
	CreatedAt string  `json:"created_at" db:"created_at"`
	UpdatedAt string  `json:"updated_at" db:"updated_at"`
}
