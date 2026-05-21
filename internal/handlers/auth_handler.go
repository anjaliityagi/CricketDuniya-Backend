package handlers

import (
	"CricketDuniya-Backend/internal/dto"
	"CricketDuniya-Backend/internal/repositories"
	"CricketDuniya-Backend/internal/services"

	//"CricketDuniya-Backend/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Signup(c *gin.Context) {
	var req dto.SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "Please provide valid signup details", err)
		return
	}

	user, err := services.Signup(req)
	if err != nil {
		badRequest(c, "Unable to create account with provided details", err)
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
		badRequest(c, "Please provide valid login details", err)
		return
	}

	token, err := services.Login(req)
	if err != nil {
		unauthorized(c, "Invalid phone number or password", err)
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
		internalServerError(c, "Unable to logout right now. Please try again", err)
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
