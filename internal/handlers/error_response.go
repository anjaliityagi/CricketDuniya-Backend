package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func badRequest(c *gin.Context, userMessage string, err error) {
	if err != nil {
		log.Printf("bad request at %s %s: %v", c.Request.Method, c.Request.URL.Path, err)
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"success": false,
		"message": userMessage,
	})
}

func unauthorized(c *gin.Context, userMessage string, err error) {
	if err != nil {
		log.Printf("unauthorized at %s %s: %v", c.Request.Method, c.Request.URL.Path, err)
	}
	c.JSON(http.StatusUnauthorized, gin.H{
		"success": false,
		"message": userMessage,
	})
}

func notFound(c *gin.Context, userMessage string, err error) {
	if err != nil {
		log.Printf("not found at %s %s: %v", c.Request.Method, c.Request.URL.Path, err)
	}
	c.JSON(http.StatusNotFound, gin.H{
		"success": false,
		"message": userMessage,
	})
}

func internalServerError(c *gin.Context, userMessage string, err error) {
	if err != nil {
		log.Printf("internal error at %s %s: %v", c.Request.Method, c.Request.URL.Path, err)
	}
	c.JSON(http.StatusInternalServerError, gin.H{
		"success": false,
		"message": userMessage,
	})
}
