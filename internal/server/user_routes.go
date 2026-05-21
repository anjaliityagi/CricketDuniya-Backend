package server

import (
	"CricketDuniya-Backend/internal/handlers"
	"CricketDuniya-Backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func registerUserRoutes(rg *gin.RouterGroup) {
	users := rg.Group("/")
	users.Use(middleware.AuthMiddleware())

	users.GET("/users", handlers.GetUsers)
}
