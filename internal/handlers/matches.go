package handlers

import (
	"CricketDuniya-Backend/internal/dto"
	"CricketDuniya-Backend/internal/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateMatch(c *gin.Context) {
	var req dto.CreateMatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "Please provide valid match details", err)
		return
	}
	userID := c.GetString("user_id")
	//fmt.Println(userID + "kdcksjdnvk")
	match, err := services.CreateMatch(req, userID)
	fmt.Println(userID)
	if err != nil {
		internalServerError(c, "Unable to create match right now. Please try again", err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "match created successfully",
		"data":    match,
	})
}
func GetAllMatches(c *gin.Context) {
	var query dto.GetMatchesQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		badRequest(c, "Please provide valid query parameters", err)
		return
	}
	matches, err := services.GetAllMatches(query)
	if err != nil {
		internalServerError(c, "Unable to fetch matches right now. Please try again", err)

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"matches": matches,
	})
}

func GetMatchByID(c *gin.Context) {

	id := c.Param("id")

	match, err := services.GetMatchDetailByID(id)
	if err != nil {
		notFound(c, "Match not found", err)
		fmt.Println(err)
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"data":    match,
	})
}

func CompleteMatch(c *gin.Context) {
	matchID := c.Param("id")
	if matchID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "match id is required"})
		return
	}

	var req dto.CompleteMatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "Please provide valid match completion details", err)
		return
	}

	winnerTeamID := req.WinnerMatchTeamID
	if winnerTeamID == "" {
		winnerTeamID = req.WinnerTeamID
	}
	if winnerTeamID == "" {
		badRequest(c, "winner_match_team_id or winner_team_id is required", nil)
		return
	}

	if err := services.CompleteMatch(matchID, winnerTeamID); err != nil {
		internalServerError(c, "Unable to complete match right now. Please try again", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "match completed successfully",
	})
}

func GetMatchInnings(c *gin.Context) {
	matchID := c.Param("id")
	if matchID == "" {
		badRequest(c, "match id is required", nil)
		return
	}

	innings, err := services.GetMatchInnings(matchID)
	if err != nil {
		internalServerError(c, "Unable to fetch innings right now. Please try again", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"innings": innings,
	})
}

func GetMatchScorecard(c *gin.Context) {
	matchID := c.Param("id")
	if matchID == "" {
		badRequest(c, "match id is required", nil)
		return
	}

	scorecard, err := services.GetMatchScorecard(matchID)
	if err != nil {
		internalServerError(c, "Unable to fetch scorecard right now. Please try again", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    scorecard,
	})
}

func StartMatch(c *gin.Context) {
	matchID := c.Param("id")
	if matchID == "" {
		badRequest(c, "match id is required", nil)
		return
	}

	innings, err := services.StartMatch(matchID)
	if err != nil {
		internalServerError(c, "Unable to start match right now. Please try again", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Match started successfully",
		"innings": innings,
	})
}

func GetMatchSquad(c *gin.Context) {
	matchID := c.Param("id")
	if matchID == "" {
		badRequest(c, "match id is required", nil)
		return
	}

	squad, err := services.GetMatchSquad(matchID)
	if err != nil {
		internalServerError(c, "Unable to fetch match squad right now. Please try again", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"squad":   squad,
	})
}

func UpdateMatchLineup(c *gin.Context) {
	matchID := c.Param("id")
	if matchID == "" {
		badRequest(c, "match id is required", nil)
		return
	}

	var req dto.UpdateMatchLineupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "Please provide valid lineup details", err)
		return
	}

	if err := services.UpdateMatchLineup(matchID, req); err != nil {
		internalServerError(c, "Unable to update lineup right now. Please try again", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Lineup updated successfully",
	})
}

func SetFirstPickTeam(c *gin.Context) {
	matchID := c.Param("id")
	if matchID == "" {
		badRequest(c, "match id is required", nil)
		return
	}

	var req dto.SetFirstPickRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "Please provide valid first pick details", err)
		return
	}

	match, err := services.SetFirstPickTeam(matchID, req.FirstPickTeamID)
	if err != nil {
		if err.Error() == "first pick team must be one of the match teams" {
			badRequest(c, "First pick team must be one of the match teams", err)
			return
		}
		internalServerError(c, "Unable to save first pick right now. Please try again", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "First pick team saved successfully",
		"data":    match,
	})
}
