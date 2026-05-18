package server

import (
	"CricketDuniya-Backend/internal/handlers"

	"github.com/gin-gonic/gin"
)

func MatchRoutes(r *gin.RouterGroup) {

	r.GET("/matches", handlers.GetAllMatches)
	r.GET("/matches/:id", handlers.GetMatchByID)
}
