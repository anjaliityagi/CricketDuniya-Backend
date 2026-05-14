package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	Router *gin.Engine
}

func SetupRoutes() *Server {

	router := gin.Default()

	v1 := router.Group("/v1")
	{

		v1.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status": "Server is Running",
			})
		})

	}

	return &Server{
		Router: router,
	}
}
