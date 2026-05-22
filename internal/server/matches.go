package server

import (
	"CricketDuniya-Backend/internal/handlers"

	"github.com/gin-gonic/gin"
)

func MatchRoutes(r *gin.RouterGroup) {

	r.GET("/matches", handlers.GetAllMatches)
	r.GET("/matches/:id", handlers.GetMatchByID)
	r.GET("/matches/:id/innings", handlers.GetMatchInnings)
	r.GET("/matches/:id/scorecard", handlers.GetMatchScorecard)
	r.GET("/matches/:id/squad", handlers.GetMatchSquad)
}
