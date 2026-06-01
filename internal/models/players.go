package models

import "time"

type MatchPlayer struct {
	ID       string  `db:"id" json:"id"`
	PlayerID *string `db:"player_id" json:"player_id,omitempty"`

	TeamID    *string    `db:"team_id" json:"team_id"`
	IsHost    bool       `db:"is_host" json:"is_host"`
	IsCaptain bool       `db:"is_captain" json:"is_captain"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	RemovedAt *time.Time `db:"removed_at" json:"removed_at,omitempty"`
}

type TeamPlayer struct {
	ID       string  `db:"id"`
	TeamID   string  `db:"team_id"`
	PlayerID *string `db:"player_id"`
}
