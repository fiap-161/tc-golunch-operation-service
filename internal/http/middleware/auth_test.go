package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	authcontroller "github.com/fiap-161/tc-golunch-operation-service/internal/auth/controller"
	"github.com/fiap-161/tc-golunch-operation-service/internal/auth/external"

	"github.com/gin-gonic/gin"
)

func TestAuthMiddleware(t *testing.T) {
	jwtGateway := external.NewJWTService("secret", time.Minute*5)
	controller := authcontroller.New(jwtGateway)
	validToken, err := controller.GenerateToken("user123", "admin", nil)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	tests := []struct {
		name               string
		authHeader         string
		expectedStatusCode int
	}{
		{
			name:               "missing Authorization header",
			authHeader:         "",
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name:               "invalid format in Authorization header",
			authHeader:         "InvalidHeader",
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name:               "invalid token value",
			authHeader:         "Bearer invalid.token.value",
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name:               "valid token",
			authHeader:         "Bearer " + validToken,
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			router := gin.New()
			router.Use(AuthMiddleware(controller))
			router.GET("/protected", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "OK"})
			})

			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			if resp.Code != tt.expectedStatusCode {
				t.Errorf("expected status %d, got %d", tt.expectedStatusCode, resp.Code)
			}
		})
	}
}
