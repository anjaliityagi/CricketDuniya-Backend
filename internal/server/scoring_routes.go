package server

import (
	"CricketDuniya-Backend/internal/handlers"

	"github.com/gin-gonic/gin"
)

func ScoringRoutes(r *gin.RouterGroup) {

	//h := handlers.NewBallHandler()
	scoring := r.Group("/")
	{
		scoring.POST("/ball", handlers.AddBall)
		scoring.PATCH("/innings/:innings_id/state", handlers.OverrideInningsState)
		scoring.POST("/innings/:innings_id/undo-ball", handlers.UndoLastBall)
	}
}
