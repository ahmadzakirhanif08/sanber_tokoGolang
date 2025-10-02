package middlewares

import (
	"os"
	"net/http"
	"github.com/gin-gonic/gin"
)


func BasicAuthMiddleware() gin.HandlerFunc {
	authUsername := os.Getenv("BASIC_AUTH_USER")
	authPassword := os.Getenv("BASIC_AUTH_PASSWORD")
	
	if authUsername == "" || authPassword == "" {
		return func(c *gin.Context) {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Basic Auth credentials not set in environment"})
		}
	}

	return func(c *gin.Context) {
		user, password, hasAuth := c.Request.BasicAuth()

		if !hasAuth || user != authUsername || password != authPassword {
			c.Header("WWW-Authenticate", "Basic realm=Restricted")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Invalid credentials"})
			return
		}
		c.Next()
	}
}