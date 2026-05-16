package dto

type CreatePlayerRequest struct {
	MatchID    string  `json:"match_id" binding:"required"`
	PlayerName string  `json:"player_name" binding:"required"`
	Phone      *string `json:"phone"`
	TeamSide   *string `json:"team_side"`
}
