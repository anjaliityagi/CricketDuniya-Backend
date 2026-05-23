package dto

type CreateMatchWithTeamsRequest struct {
	TeamAName       string `json:"team_a_name" binding:"required"`
	TeamBName       string `json:"team_b_name" binding:"required"`
	Location        string `json:"location"`
	MatchDate       string `json:"match_date"`
	OversPerInnings int    `json:"overs_per_innings"`
}
