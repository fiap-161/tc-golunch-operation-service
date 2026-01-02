package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// ServiceAuthMiddleware validates service-to-service authentication
func ServiceAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip authentication for health checks and public endpoints
		if c.Request.URL.Path == "/ping" || c.Request.URL.Path == "/health" {
			c.Next()
			return
		}

		// Method 1: API Key Authentication (preferred for service-to-service)
		serviceName := c.GetHeader("X-Service-Name")
		serviceKey := c.GetHeader("X-Service-Key")

		if serviceName != "" && serviceKey != "" {
			if validateServiceAPIKey(serviceName, serviceKey) {
				c.Set("authenticated_service", serviceName)
				c.Next()
				return
			}
		}

		// Method 2: JWT Token Authentication (for user requests)
		authHeader := c.GetHeader("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			// Let the existing auth middleware handle JWT tokens
			c.Next()
			return
		}

		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Invalid service credentials"})
	}
}

// validateServiceAPIKey validates API key for service authentication
func validateServiceAPIKey(serviceName, apiKey string) bool {
	var expectedKey string

	switch serviceName {
	case "core-service":
		expectedKey = os.Getenv("CORE_SERVICE_API_KEY")
	case "payment-service":
		expectedKey = os.Getenv("PAYMENT_SERVICE_API_KEY")
	case "operation-service":
		expectedKey = os.Getenv("OPERATION_SERVICE_API_KEY")
	default:
		return false
	}

	if expectedKey == "" {
		return false
	}

	// Use constant-time comparison to prevent timing attacks
	return constantTimeEquals(apiKey, expectedKey)
}

// constantTimeEquals performs constant-time string comparison
func constantTimeEquals(a, b string) bool {
	if len(a) != len(b) {
		return false
	}

	result := 0
	for i := 0; i < len(a); i++ {
		result |= int(a[i]) ^ int(b[i])
	}

	return result == 0
}
