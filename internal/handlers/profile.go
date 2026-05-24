package handlers

import (
	"CricketDuniya-Backend/internal/dto"
	"CricketDuniya-Backend/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUserProfile(c *gin.Context) {
	userID := c.GetString("user_id")

	profile, err := services.GetUserProfile(userID)
	if err != nil {
		internalServerError(c, "Unable to fetch profile right now. Please try again", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "profile fetched successfully",
		"data":    profile,
	})
}

func GetPublicUserProfile(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		badRequest(c, "Please provide a valid player id", nil)
		return
	}

	profile, err := services.GetPublicUserProfile(userID)
	if err != nil {
		internalServerError(c, "Unable to fetch player profile right now. Please try again", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "player profile fetched successfully",
		"data":    profile,
	})
}

func UpdateUserProfile(c *gin.Context) {
	var req dto.UpdateProfileRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "Please provide valid profile details", err)
		return
	}

	userID := c.GetString("user_id")

	profile, err := services.UpdateUserProfile(userID, req)
	if err != nil {
		internalServerError(c, "Unable to update profile right now. Please try again", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "profile updated successfully",
		"data":    profile,
	})
}
