package server

import (
	"CricketDuniya-Backend/internal/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterWinProbabilityRoutes(r *gin.RouterGroup) {
	winHandler := handlers.NewWinProbabilityHandler()
	matches := r.Group("/matches")
	{
		matches.GET("/:id/win-probability", winHandler.GetWinProbability)
	}
}
