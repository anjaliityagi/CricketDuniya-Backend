package handlers

import (
	"CricketDuniya-Backend/internal/dto"
	"CricketDuniya-Backend/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateTeam(c *gin.Context) {

	var req dto.CreateTeamRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDAny, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	userID := userIDAny.(string)

	team, err := services.CreateTeam(req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, team)
}

func GetTeam(c *gin.Context) {

	id := c.Param("id")

	team, err := services.GetTeamById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "team not found"})
		return
	}

	c.JSON(http.StatusOK, team)
}

func GetTeams(c *gin.Context) {

	teams, err := services.GetTeams()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, teams)
}
