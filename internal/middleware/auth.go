package middleware

import (
	"CricketDuniya-Backend/internal/repositories"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(401, gin.H{"error": "missing token"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(401, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		claims := token.Claims.(jwt.MapClaims)

		userID := fmt.Sprintf("%v", claims["user_id"])

		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token claims",
			})
			c.Abort()
			return
		}

		sessionID := claims["session_id"].(string)

		// 🔥 SESSION CHECK
		active, _ := repositories.IsSessionActive(sessionID)
		if !active {
			c.JSON(401, gin.H{"error": "session expired"})
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Set("session_id", sessionID)

		c.Next()
	}
}
