package dto

type TossRequest struct {
	MatchID          string `json:"match_id" binding:"required"`
	TossWinnerTeamID string `json:"toss_winner_team_id"`
	Decision         string `json:"decision" binding:"required,oneof=bat bowl"`
}
