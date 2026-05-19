package dto

import "time"

type CreateMatchRequest struct {
	Venue        string `json:"venue"`
	MatchDate    string `json:"match_date" binding:"required"`
	OversPerSide int    `json:"overs_per_side" binding:"required,min=1"`
	TeamAID      string `json:"team_a_id" binding:"required"`
	TeamBID      string `json:"team_b_id" binding:"required"`
}

type GetMatchesQuery struct {
	Status     string `form:"status"`
	TeamID     string `form:"team_id"`
	HostUserID string `form:"host_user_id"`
	Search     string `form:"search"`
}

type MatchResponse struct {
	ID string `db:"id" json:"id"`

	TeamAID   *string `db:"team_a_id" json:"team_a_id"`
	TeamAName *string `db:"team_a_name" json:"team_a_name"`

	TeamBID   *string `db:"team_b_id" json:"team_b_id"`
	TeamBName *string `db:"team_b_name" json:"team_b_name"`

	Venue *string `db:"venue" json:"venue"`

	Status string `db:"status" json:"status"`

	MatchDate *time.Time `db:"match_date" json:"match_date"`

	OversPerSide int `db:"overs_per_side" json:"overs_per_side"`

	TossDecision *string `db:"toss_decision" json:"toss_decision"`

	WinnerTeamID *string `db:"winner_team_id" json:"winner_team_id"`
}
