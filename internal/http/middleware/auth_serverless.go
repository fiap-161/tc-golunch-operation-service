package middleware

import (
	"net/http"
	"strings"

	"github.com/fiap-161/tc-golunch-operation-service/internal/shared/gateway"
	"github.com/gin-gonic/gin"
)

// ServerlessAuthMiddleware validates JWT tokens via serverless auth
// Following the exact same pattern as tc-golunch-api monolith
func ServerlessAuthMiddleware(authGateway *gateway.ServerlessAuthGateway) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			return
		}

		tokenString := parts[1]

		claims, err := authGateway.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// Set context values exactly like monolith
		c.Set("user_id", claims.UserID)
		c.Set("user_type", claims.UserType)
		c.Set("claims", claims)

		c.Next()
	}
}

// ServerlessAdminOnly middleware to restrict access to admin users only
// Following the same pattern as tc-golunch-api monolith
func ServerlessAdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		userType, exists := c.Get("user_type")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User type not found in context"})
			return
		}

		if userType != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			return
		}

		c.Next()
	}
}