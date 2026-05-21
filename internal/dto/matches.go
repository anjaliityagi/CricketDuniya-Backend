package dto

import "time"

type CreateMatchRequest struct {
	TeamAID         string `json:"team_a_id"`
	TeamBID         string `json:"team_b_id"`
	Location        string `json:"location"`
	MatchDate       string `json:"match_date" binding:"required"`
	OversPerInnings int    `json:"overs_per_innings" binding:"required,min=1"`
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

	Location *string `db:"location" json:"location"`

	Status string `db:"status" json:"status"`

	MatchDate *time.Time `db:"match_date" json:"match_date"`

	OversPerInnings int `db:"overs_per_innings" json:"overs_per_innings"`

	TossDecision *string `db:"toss_decision" json:"toss_decision"`

	WinnerTeamID *string `db:"winner_match_team_id" json:"winner_match_team_id"`
}

type CompleteMatchRequest struct {
	WinnerMatchTeamID string `json:"winner_match_team_id" binding:"required"`
}
