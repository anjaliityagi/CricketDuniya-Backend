package handlers

import (
	"CricketDuniya-Backend/internal/database"
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

	tx, err := database.DB.Beginx()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to start transaction",
		})
		return
	}

	defer tx.Rollback()

	teamA, err := services.CreateTeam(
		dto.CreateTeamRequest{
			Name: req.TeamAName,
		},
		userID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create team A",
		})
		return
	}

	teamB, err := services.CreateTeam(
		dto.CreateTeamRequest{
			Name: req.TeamBName,
		},
		userID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create team B",
		})
		return
	}

	match, err := services.CreateMatch(
		dto.CreateMatchRequest{
			TeamAID:         *teamA.ID,
			TeamBID:         *teamB.ID,
			Location:        req.Location,
			MatchDate:       req.MatchDate,
			OversPerInnings: req.OversPerInnings,
		},
		userID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create match",
		})
		return
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to commit transaction",
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
