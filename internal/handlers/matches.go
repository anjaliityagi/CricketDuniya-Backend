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

		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	userID := c.GetString("user_id")
	//fmt.Println(userID + "kdcksjdnvk")
	match, err := services.CreateMatch(req, userID)
	fmt.Println(userID)
	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "match created successfully",
		"data":    match,
	})
}
