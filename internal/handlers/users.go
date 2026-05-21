package handlers

import (
	"CricketDuniya-Backend/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context) {
	users, err := services.GetAllUsers()
	if err != nil {
		internalServerError(c, "Unable to fetch users right now. Please try again", err)
		return
	}

	c.JSON(http.StatusOK, users)
}
