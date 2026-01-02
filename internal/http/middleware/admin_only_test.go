package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestAdminOnly(t *testing.T) {
	tests := []struct {
		name               string
		userType           any
		setUserType        bool
		expectedStatusCode int
	}{
		{
			name:               "missing user_type in context",
			setUserType:        false,
			expectedStatusCode: http.StatusForbidden,
		},
		{
			name:               "non-string user_type",
			setUserType:        true,
			userType:           123,
			expectedStatusCode: http.StatusForbidden,
		},
		{
			name:               "non-admin user_type",
			setUserType:        true,
			userType:           "customer",
			expectedStatusCode: http.StatusForbidden,
		},
		{
			name:               "admin user_type",
			setUserType:        true,
			userType:           "admin",
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			router := gin.New()

			router.Use(func(c *gin.Context) {
				if tt.setUserType {
					c.Set("user_type", tt.userType)
				}
			}, AdminOnly())

			router.GET("/admin", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "Access granted"})
			})

			req := httptest.NewRequest(http.MethodGet, "/admin", nil)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			if resp.Code != tt.expectedStatusCode {
				t.Errorf("expected status %d, got %d", tt.expectedStatusCode, resp.Code)
			}
		})
	}
}
