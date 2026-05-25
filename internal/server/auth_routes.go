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
		//auth.POST("/forgetPassword", handlers.ForgetPassword)

		auth.POST("/forgot-password", handlers.ForgotPassword)
		auth.POST("/verify-otp", handlers.VerifyOTPAndResetPassword)

		protected := rg.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.POST("/logout", handlers.Logout)

			protected.GET("/profile", handlers.GetUserProfile)

			protected.PATCH("/profile", handlers.UpdateUserProfile)

			//protected.POST("/matches", handlers.CreateMatch)

			protected.POST("/matches", handlers.CreateMatchWithTeams)

			protected.PATCH("/matches/:id/first-pick", handlers.SetFirstPickTeam)

			protected.PATCH("/matches/:id/start", handlers.StartMatch)

			protected.PATCH("/matches/:id/lineup", handlers.UpdateMatchLineup)
			protected.PATCH("/matches/:id/complete", handlers.CompleteMatch)
			//
			//protected.POST("/teams/:id/players", handlers.CreatePlayer)
			protected.DELETE("/players/:id", handlers.DeletePlayer)
			//protected.GET("/matches", handlers.GetAllMatches)

			//protected.POST("/matches/:id/toss", handlers.DoToss)
			protected.POST("/toss", handlers.TossHandler)

			//protected.POST("/api/v1/ball/update", handlers.UpdateBall)

			protected.POST("/team-players", handlers.AddPlayersToTeams)
			protected.PUT("players/:player_id/assign-captain", handlers.AssignCaptain)
			protected.PUT("players/:player_id/assign-wicketkeeper", handlers.AssignWicketKeeper)

		}
	}
}
