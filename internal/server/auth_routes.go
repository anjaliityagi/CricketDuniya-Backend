package server

import (
	"CricketDuniya-Backend/internal/handlers"
	"CricketDuniya-Backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func registerAuthRoutes(rg *gin.RouterGroup) {
	rg.GET("/users/:id/profile", handlers.GetPublicUserProfile)

	rg.GET("/players", handlers.GetPlayersDirectory)

	auth := rg.Group("/auth")
	{
		auth.POST("/signup", handlers.Signup)
		auth.POST("/login", handlers.Login)

		auth.POST("/forgot-password", handlers.ForgotPassword)
		auth.POST("/verify-otp", handlers.VerifyOTPAndResetPassword)

		protected := rg.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.POST("/logout", handlers.Logout)

			protected.GET("/profile", handlers.GetUserProfile)

			protected.PATCH("/profile", handlers.UpdateUserProfile)

			protected.POST("/matches", handlers.CreateMatchWithTeams)

			protected.PATCH("/matches/:id/first-pick", handlers.SetFirstPickTeam)

			protected.PATCH("/matches/:id/start", handlers.StartMatch)

			protected.PATCH("/matches/:id/lineup", handlers.UpdateMatchLineup)
			protected.PATCH("/matches/:id/complete", handlers.CompleteMatch)
			//protected.DELETE("/players/:id", handlers.DeletePlayer)

			protected.POST("/toss", handlers.TossHandler)
			//
			//protected.POST("/team-players", handlers.AddPlayersToTeams)
			//protected.PUT("players/:player_id/assign-captain", handlers.AssignCaptain)
			//protected.PUT("players/:player_id/assign-umpire", handlers.AssignUmpire)

		}
	}
}
