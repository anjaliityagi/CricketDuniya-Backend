package handlers

import (
	"CricketDuniya-Backend/internal/dto"
	"CricketDuniya-Backend/internal/services"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreatePlayer(c *gin.Context) {

	var req dto.CreatePlayerRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	userID := c.GetString("user_id")

	player, err := services.CreatePlayer(req, userID)
	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "player added successfully",
		"data":    player,
	})
}

func DeletePlayer(c *gin.Context) {
	playerID := c.Param("id")
	if playerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "player id is required",
		})
		return
	}

	userID := c.GetString("user_id")

	err := services.DeletePlayer(playerID, userID)
	if err != nil {
		if errors.Is(err, services.ErrPlayerNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "player removed successfully",
	})
}
