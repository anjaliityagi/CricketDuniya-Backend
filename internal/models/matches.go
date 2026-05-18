package models

import "time"

type Match struct {
	ID string `db:"id" json:"id"`

	HostUserID string `db:"host_user_id" json:"host_user_id"`

	TeamAID *string `db:"team_a_id" json:"team_a_id"`
	TeamBID *string `db:"team_b_id" json:"team_b_id"`

	Venue *string `db:"venue" json:"venue,omitempty"`

	MatchDate *time.Time `db:"match_date" json:"match_date,omitempty"`

	OversPerSide int `db:"overs_per_side" json:"overs_per_side"`

	TossDecision *string `db:"toss_decision" json:"toss_decision,omitempty"`

	TossWinnerTeamID *string `db:"toss_winner_team_id" json:"toss_winner_team_id,omitempty"`

	WinnerTeamID *string `db:"winner_team_id" json:"winner_team_id,omitempty"`

	Status string `db:"status" json:"status"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
