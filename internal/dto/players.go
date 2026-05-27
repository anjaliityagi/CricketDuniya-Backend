package dto

import "time"

type CreatePlayerRequest struct {
	PlayerID     *string `json:"player_id"`
	Name         *string `json:"name"`
	PhoneNumber  *string `json:"phone_number"`
	BattingStyle *string `json:"batting_style"`
	BowlingStyle *string `json:"bowling_style"`
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
	IsUmpire       bool       `db:"is_umpire" json:"is_umpire"`
	IsPlayingXI    bool       `db:"is_playing_xi" json:"is_playing_xi"`
	BattingOrder   *int       `db:"batting_order" json:"batting_order,omitempty"`
	CreatedAt      time.Time  `db:"created_at" json:"created_at"`
	RemovedAt      *time.Time `db:"removed_at" json:"removed_at,omitempty"`
}
type AddPlayersRequest []TeamPlayersRequest

type TeamPlayersRequest struct {
	TeamID  string        `json:"team_id" binding:"required"`
	Players []PlayerInput `json:"players" binding:"required"`
}

type PlayerInput struct {
	PlayerID    *string `json:"player_id,omitempty"`
	Name        string  `json:"name,omitempty"`
	PhoneNumber string  `json:"phone_number,omitempty"`
}

type AddTeamPlayerRequest struct {
	PlayerID    *string `json:"player_id,omitempty"`
	Name        string  `json:"name,omitempty"`
	PhoneNumber string  `json:"phone_number,omitempty"`
}

type AssignCaptainRequest struct {
	TeamID string `json:"team_id" binding:"required"`
}

type AssignUmpireRequest struct {
	TeamID string `json:"team_id" binding:"required"`
}

type PlayerDirectoryItem struct {
	ID            string  `db:"id" json:"id"`
	Name          string  `db:"name" json:"name"`
	PhoneNumber   *string `db:"phone_number" json:"phone_number,omitempty"`
	BattingStyle  *string `db:"batting_style" json:"batting_style,omitempty"`
	BowlingStyle  *string `db:"bowling_style" json:"bowling_style,omitempty"`
	MatchesPlayed int     `db:"matches_played" json:"matches_played"`
	Points        int     `db:"points" json:"points"`
}
