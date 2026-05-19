package dto

type CreateTeamRequest struct {
	Name      string `json:"name" binding:"required"`
	CaptionID string `json:"caption_id" binding:"required"`
}
