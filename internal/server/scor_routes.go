package server

import (
	"CricketDuniya-Backend/internal/handlers"

	"github.com/gin-gonic/gin"
)

func ScoringRoutes(r *gin.RouterGroup) {

	scoring := r.Group("/")
	{
		scoring.POST("/ball", handlers.NewBallHandler().AddBall)
	}
}
