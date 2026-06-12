package server

import (
	"CricketDuniya-Backend/internal/handlers"

	"github.com/gin-gonic/gin"
)

func registerTeamRoutes(rg *gin.RouterGroup) {

	teams := rg.Group("/teams")
	//teams.Use(middleware.AuthMiddleware())

	teams.POST("", handlers.CreateTeam)
	teams.GET("/:id", handlers.GetTeam)
	teams.GET("/:id/players", handlers.GetTeamPlayers)
	teams.GET("", handlers.GetTeams)
}
