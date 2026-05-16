package models

import "time"

type MatchPlayer struct {
	ID             string    `db:"id" json:"id"`
	MatchID        string    `db:"match_id" json:"match_id"`
	UserID         *string   `db:"user_id" json:"user_id,omitempty"`
	PlayerName     string    `db:"player_name" json:"player_name"`
	Phone          *string   `db:"phone" json:"phone,omitempty"`
	TeamSide       *string   `db:"team_side" json:"team_side,omitempty"`
	IsHost         bool      `db:"is_host" json:"is_host"`
	IsCaptain      bool      `db:"is_captain" json:"is_captain"`
	IsWicketkeeper bool      `db:"is_wicketkeeper" json:"is_wicketkeeper"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
}
