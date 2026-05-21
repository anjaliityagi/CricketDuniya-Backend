package handlers

import (
	"CricketDuniya-Backend/internal/dto"
	"CricketDuniya-Backend/internal/repositories"
	"CricketDuniya-Backend/internal/services/matchscoring"
	"CricketDuniya-Backend/internal/services/scoring"

	"github.com/gin-gonic/gin"
)

type BallHandler struct {
	fantasyEngine *scoring.Engine
	matchEngine   *matchscoring.Engine
}

func NewBallHandler() *BallHandler {
	return &BallHandler{
		fantasyEngine: scoring.NewEngine(),
		matchEngine:   matchscoring.NewEngine(),
	}
}

func (h *BallHandler) AddBall(c *gin.Context) {
	var req dto.BallRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "Please provide valid ball data", err)
		return
	}

	matchUpdate, err := h.matchEngine.Process(req)
	if err != nil {
		internalServerError(c, "Unable to process ball right now. Please try again", err)
		return
	}

	bat, bowl, field, err := h.fantasyEngine.Process(req)
	if err != nil {
		internalServerError(c, "Unable to calculate points right now. Please try again", err)
		return
	}

	ballEventID, err := repositories.SaveBall(req)
	if err != nil {
		internalServerError(c, "Unable to save ball right now. Please try again", err)
		return
	}

	if err := scoring.PersistBallFantasy(req, ballEventID, bat, bowl, field); err != nil {
		internalServerError(c, "Unable to save fantasy points right now. Please try again", err)
		return
	}

	c.JSON(200, gin.H{
		"message":         "ball saved",
		"match_score":     matchUpdate,
		"batting_points":  bat,
		"bowling_points":  bowl,
		"fielding_points": field,
	})
}
