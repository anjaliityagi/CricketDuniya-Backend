package dto

import "time"

type CreatePlayerRequest struct {
	PLayerID *string `json:"player_id"`
}

type TeamPlayerResponse struct {
	ID             string     `db:"id" json:"id"`
	TeamID         string     `db:"team_id" json:"team_id"`
	PlayerID       *string    `db:"player_id" json:"player_id"`
	Name           string     `db:"name" json:"name"`
	PhoneNumber    *string    `db:"phone_number" json:"phone_number,omitempty"`
	BattingStyle   *string    `db:"batting_style" json:"batting_style,omitempty"`
	BowlingStyle   *string    `db:"bowling_style" json:"bowling_style,omitempty"`
	IsCaptain      bool       `db:"is_captain" json:"is_captain"`
	IsWicketKeeper bool       `db:"is_wicket_keeper" json:"is_wicket_keeper"`
	IsPlayingXI    bool       `db:"is_playing_xi" json:"is_playing_xi"`
	IsSubstitute   bool       `db:"is_substitute" json:"is_substitute"`
	BattingOrder   *int       `db:"batting_order" json:"batting_order,omitempty"`
	CreatedAt      time.Time  `db:"created_at" json:"created_at"`
	RemovedAt      *time.Time `db:"removed_at" json:"removed_at,omitempty"`
}
