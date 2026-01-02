package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		userType, exists := c.Get("user_type")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "User type not found"})
			return
		}

		userTypeStr, ok := userType.(string)
		if !ok || strings.ToLower(userTypeStr) != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			return
		}

		c.Next()
	}
}
