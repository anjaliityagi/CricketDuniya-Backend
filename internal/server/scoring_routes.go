package server

import (
	"CricketDuniya-Backend/internal/handlers"

	"github.com/gin-gonic/gin"
)

func ScoringRoutes(r *gin.RouterGroup) {

	h := handlers.NewBallHandler()
	scoring := r.Group("/")
	{
		scoring.POST("/ball", h.AddBall)
		scoring.PATCH("/innings/:innings_id/state", h.OverrideInningsState)
		scoring.POST("/innings/:innings_id/undo-ball", h.UndoLastBall)
	}
}
