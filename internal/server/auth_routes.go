package server

import (
	"CricketDuniya-Backend/internal/handlers"

	"github.com/gin-gonic/gin"
)

func registerAuthRoutes(rg *gin.RouterGroup) {

	auth := rg.Group("/auth")
	{
		auth.POST("/signup", handlers.Signup)
		auth.POST("/login", handlers.Login)
	}
}
