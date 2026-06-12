package server

import (
	"CricketDuniya-Backend/internal/handlers"
	"CricketDuniya-Backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func registerPlayerRoutes(rg *gin.RouterGroup) {
	//rg.GET("/players", handlers.GetPlayersDirectory)

	players := rg.Group("/")
	players.Use(middleware.AuthMiddleware())
	{
		players.DELETE("/players/:id", handlers.DeletePlayer)
		players.POST("/team-players", handlers.AddPlayersToTeams)
		players.PUT("players/:player_id/assign-captain", handlers.AssignCaptain)
		players.PUT("players/:player_id/assign-umpire", handlers.AssignUmpire)

	}
}
