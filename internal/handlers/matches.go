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
func GetAllMatches(c *gin.Context) {
	var query dto.GetMatchesQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	matches, err := services.GetAllMatches(query)
	if err != nil {
		//fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"matches": matches,
	})
}

func GetMatchByID(c *gin.Context) {

	id := c.Param("id")

	match, err := services.GetMatchByID(id)
	if err != nil {
		c.JSON(404, gin.H{
			"success": false,
			"message": "match not found",
		})
		fmt.Println(err)
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"data":    match,
	})
}
