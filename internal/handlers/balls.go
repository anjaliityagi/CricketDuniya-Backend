package handlers

import (
	"CricketDuniya-Backend/internal/dto"
	"CricketDuniya-Backend/internal/services/matchscoring"
	"strings"

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
		if isScoringClientError(err) {
			badRequest(c, err.Error(), err)
			return
		}
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
		if isScoringClientError(err) {
			badRequest(c, err.Error(), err)
			return
		}
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

func isScoringClientError(err error) bool {
	if err == nil {
		return false
	}

	msg := strings.ToLower(err.Error())
	clientFragments := []string{
		"same bowler cannot bowl consecutive overs",
		"next_batter_id is required",
		"innings is not live",
		"innings already completed by overs",
		"required for first ball",
		"active striker, non_striker and bowler are required",
		"striker_id, non_striker_id and bowler_id are required",
		"strikers must belong to batting team",
		"bowler must belong to bowling team",
		"striker and non_striker must be different players",
	}
	for _, fragment := range clientFragments {
		if strings.Contains(msg, fragment) {
			return true
		}
	}
	return false
}
