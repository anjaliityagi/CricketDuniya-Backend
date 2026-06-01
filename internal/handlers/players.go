package handlers

import (
	"CricketDuniya-Backend/internal/dto"
	"CricketDuniya-Backend/internal/services"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetPlayersDirectory(c *gin.Context) {
	search := c.Query("search")
	limit := 10

	if c.Query("limit") != "" {
		parsedLimit, err := strconv.Atoi(c.Query("limit"))
		if err != nil {
			badRequest(c, "Please provide a valid players limit", err)
			return
		}
		limit = parsedLimit
	}

	players, err := services.GetPlayersDirectory(search, limit)
	if err != nil {
		internalServerError(c, "Unable to fetch players right now. Please try again", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "players fetched successfully",
		"data":    players,
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

func AddPlayersToTeams(c *gin.Context) {

	var req dto.AddPlayersRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})

		return
	}

	err := services.AddPlayersToTeams(req)

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "players added successfully",
	})
}

func AssignCaptain(c *gin.Context) {

	playerID := c.Param("player_id")

	var req dto.AssignCaptainRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})

		return
	}

	err := services.AssignCaptain(
		req.TeamID,
		playerID,
	)

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "captain assigned successfully",
	})
}

func AssignUmpire(c *gin.Context) {

	playerID := c.Param("player_id")

	var req dto.AssignUmpireRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})

		return
	}

	err := services.AssignUmpire(
		req.TeamID,
		playerID,
	)

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "umpire assigned successfully",
	})
}
