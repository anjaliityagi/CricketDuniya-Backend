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

			protected.GET("/profile", handlers.GetUserProfile)

			protected.PATCH("/profile", handlers.UpdateUserProfile)

			protected.POST("/matches", handlers.CreateMatch)
			protected.PATCH("/matches/:id/start", handlers.StartMatch)
			protected.PATCH("/matches/:id/lineup", handlers.UpdateMatchLineup)
			protected.PATCH("/matches/:id/complete", handlers.CompleteMatch)

			protected.POST("/teams/:id/players", handlers.CreatePlayer)
			protected.DELETE("/players/:id", handlers.DeletePlayer)
			//protected.GET("/matches", handlers.GetAllMatches)

			//protected.POST("/matches/:id/toss", handlers.DoToss)
			protected.POST("/toss", handlers.TossHandler)

			//protected.POST("/api/v1/ball/update", handlers.UpdateBall)
		}
	}
}
