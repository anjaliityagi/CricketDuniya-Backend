package handlers

import (
	"CricketDuniya-Backend/internal/dto"
	"CricketDuniya-Backend/internal/repositories"
	"CricketDuniya-Backend/internal/services/scoring"

	"github.com/gin-gonic/gin"
)

type BallHandler struct {
	engine *scoring.Engine
}

func NewBallHandler() *BallHandler {
	return &BallHandler{
		engine: scoring.NewEngine(),
	}
}

func (h *BallHandler) AddBall(c *gin.Context) {
	var req dto.BallRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// save ball
	if err := repositories.SaveBall(req); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// scoring
	bat, bowl, field := h.engine.Process(req)

	c.JSON(200, gin.H{
		"message":         "ball saved",
		"batting_points":  bat,
		"bowling_points":  bowl,
		"fielding_points": field,
	})
}
