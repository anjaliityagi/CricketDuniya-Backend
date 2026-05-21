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

	err := services.PerformToss(req)

	if err != nil {
		internalServerError(c, "Unable to complete toss right now. Please try again", err)

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Toss completed successfully",
	})
}
