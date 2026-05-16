package server

import (
	"CricketDuniya-Backend/internal/handlers"
	"CricketDuniya-Backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func registerAuthRoutes(rg *gin.RouterGroup) {

	auth := rg.Group("/auth")
	{
		auth.POST("/signup", handlers.Signup)
		auth.POST("/login", handlers.Login)
		//auth.POST("/forgetPassword", handlers.ForgetPassword)

		protected := rg.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.POST("/logout", handlers.Logout)
			protected.POST("/matches", handlers.CreateMatch)

		}
	}
}
