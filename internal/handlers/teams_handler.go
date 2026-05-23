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
		badRequest(c, "Please provide valid team details", err)
		return
	}

	userIDAny, exists := c.Get("user_id")
	if !exists {
		unauthorized(c, "Unauthorized request", nil)
		return
	}

	userID := userIDAny.(string)

	team, err := services.CreateTeam(req, userID)
	if err != nil {
		internalServerError(c, "Unable to create team right now. Please try again", err)
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
		internalServerError(c, "Unable to fetch teams right now. Please try again", err)
		return
	}

	c.JSON(http.StatusOK, teams)
}

//func GetTeamPlayers(c *gin.Context) {
//	teamID := c.Param("id")
//	if teamID == "" {
//		c.JSON(http.StatusBadRequest, gin.H{"error": "team id is required"})
//		return
//	}
//
//	players, err := services.GetPlayersByTeam(teamID)
//	if err != nil {
//		internalServerError(c, "Unable to fetch team players right now. Please try again", err)
//		return
//	}
//
//	c.JSON(http.StatusOK, players)
//}
