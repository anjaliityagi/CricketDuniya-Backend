package handlers

import (
	"CricketDuniya-Backend/internal/dto"
	"CricketDuniya-Backend/internal/services"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreatePlayer(c *gin.Context) {

	var req dto.CreatePlayerRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "Please provide valid player details", err)
		return
	}

	teamIDStr := c.Param("id")

	teamID, err := uuid.Parse(teamIDStr)
	if err != nil {
		badRequest(c, "Invalid team id", err)
		return
	}

	player, err := services.CreatePlayer(req, teamID)
	if err != nil {
		internalServerError(c, "Unable to add player right now. Please try again", err)
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
			notFound(c, "Player not found", err)
			return
		}

		internalServerError(c, "Unable to delete player right now. Please try again", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "player removed successfully",
	})
}
