package handlers

import (
	"CricketDuniya-Backend/internal/dto"
	"CricketDuniya-Backend/internal/repositories"
	"CricketDuniya-Backend/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Signup(c *gin.Context) {
	var req dto.SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	user, err := services.Signup(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "User created successfully",
		"data":    user,
	})
}

func Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	token, err := services.Login(req)
	if err != nil {

		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Login successful",
		"token":   token,
	})
}
func Logout(c *gin.Context) {

	sessionID := c.GetString("session_id")

	err := repositories.LogoutSession(sessionID)
	if err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"message": "logout failed",
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "logged out successfully",
	})
}

//
//func forgetPassword(c *gin.Context) {
//	var req dto.LoginRequest
//	if err := c.ShouldBindJSON(&req); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{
//			"success": false,
//			"message": err.Error(),
//		})
//		return
//
//	}
//}
