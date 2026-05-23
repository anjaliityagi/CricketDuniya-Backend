package handlers

import (
	"CricketDuniya-Backend/internal/dto"
	"CricketDuniya-Backend/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateMatchWithTeams(c *gin.Context) {
	var req dto.CreateMatchWithTeamsRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	userID := c.GetString("user_id")

	teamA, teamB, match, err := services.CreateMatchWithTeams(req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create match",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "match created successfully",
		"team_a":  teamA,
		"team_b":  teamB,
		"match":   match,
	})
}
