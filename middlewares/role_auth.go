package middlewares

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func AdminAuthMiddleware() gin.HandlerFunc {
	return  func(c *gin.Context) {
		role, exists := c.Get("role")

		if !exists || role != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error":"Forbidden: Access denied. Requires Admin Role."})
			return 
		}
		c.Next()
	}
}