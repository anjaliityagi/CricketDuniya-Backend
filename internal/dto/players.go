package dto

type CreatePlayerRequest struct {
	MatchID    string  `json:"match_id" binding:"required"`
	PlayerName string  `json:"player_name" binding:"required"`
	Phone      *string `json:"phone"`
	TeamID     *string `json:"team_id"`
}
