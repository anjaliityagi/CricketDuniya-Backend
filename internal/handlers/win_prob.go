package handlers

import (
	"net/http"

	"CricketDuniya-Backend/internal/services"

	"github.com/gin-gonic/gin"
)

type WinProbabilityHandler struct{}

func NewWinProbabilityHandler() *WinProbabilityHandler {
	return &WinProbabilityHandler{}
}

func (h *WinProbabilityHandler) GetWinProbability(c *gin.Context) {

	matchID := c.Param("id")

	resp, err := services.CalculateWinProbability(matchID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}
