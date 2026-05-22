package handlers

import (
	"CricketDuniya-Backend/internal/dto"
	"CricketDuniya-Backend/internal/services/matchscoring"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BallHandler struct {
	matchEngine *matchscoring.Engine
}

func NewBallHandler() *BallHandler {
	return &BallHandler{
		matchEngine: matchscoring.NewEngine(),
	}
}

func (h *BallHandler) AddBall(c *gin.Context) {
	var req dto.BallInputRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "Please provide valid ball data", err)
		return
	}

	result, err := h.matchEngine.ProcessBall(req)
	if err != nil {
		internalServerError(c, "Unable to process ball right now. Please try again", err)
		return
	}

	c.JSON(200, gin.H{
		"message": "ball saved",
		"state":   result.State,
	})
}

func (h *BallHandler) OverrideInningsState(c *gin.Context) {
	inningsID, err := uuid.Parse(c.Param("innings_id"))
	if err != nil {
		badRequest(c, "Invalid innings id", err)
		return
	}

	var req dto.UpdateInningsStateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "Please provide valid override data", err)
		return
	}

	state, err := h.matchEngine.OverrideState(inningsID, req)
	if err != nil {
		internalServerError(c, "Unable to override innings state right now. Please try again", err)
		return
	}

	c.JSON(200, gin.H{
		"message": "innings state updated",
		"state":   state,
	})
}

func (h *BallHandler) UndoLastBall(c *gin.Context) {
	inningsID, err := uuid.Parse(c.Param("innings_id"))
	if err != nil {
		badRequest(c, "Invalid innings id", err)
		return
	}

	state, err := h.matchEngine.UndoLastBall(inningsID)
	if err != nil {
		internalServerError(c, "Unable to undo last ball right now. Please try again", err)
		return
	}

	c.JSON(200, gin.H{
		"message": "last ball undone",
		"state":   state,
	})
}
