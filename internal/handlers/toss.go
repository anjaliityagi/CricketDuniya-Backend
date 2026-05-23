package handlers

import (
	"CricketDuniya-Backend/internal/dto"
	"CricketDuniya-Backend/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func TossHandler(c *gin.Context) {

	var req dto.TossRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "Please provide valid toss details", err)

		return
	}

	innings, err := services.PerformToss(req)

	if err != nil {
		switch err.Error() {
		case "match teams are not ready",
			"invalid toss decision",
			"toss winner team id is required",
			"toss winner team is not part of this match":
			badRequest(c, err.Error(), err)
			return
		}

		internalServerError(c, "Unable to complete toss right now. Please try again", err)

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Toss completed successfully and match is now live",
		"innings": innings,
	})
}
