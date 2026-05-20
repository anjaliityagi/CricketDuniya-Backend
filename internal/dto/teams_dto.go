package dto

type CreateTeamRequest struct {
	Name string `json:"name" binding:"required"`
}
