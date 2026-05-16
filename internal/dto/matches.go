package dto

type CreateMatchRequest struct {
	Venue        string `json:"venue"`
	MatchDate    string `json:"match_date"`
	OversPerSide int    `json:"overs_per_side" binding:"required,min=1"`
}

type JoinMatchRequest struct {
	MatchID string `json:"match_id" binding:"required"`
}
