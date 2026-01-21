package gateway

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// ServerlessAuthGateway implements authentication via AWS Lambda functions
// Following the same pattern as JWTService from tc-operation-service monolith
type ServerlessAuthGateway struct {
	lambdaAuthURL  string
	serviceAuthURL string
	httpClient     *http.Client
}

// TokenRequest represents the request payload for token validation
type TokenRequest struct {
	Token string `json:"token"`
}

// TokenResponse represents the response from Lambda auth validation
type TokenResponse struct {
	Valid  bool          `json:"valid"`
	Claims *CustomClaims `json:"claims,omitempty"`
	Error  string        `json:"error,omitempty"`
}

// CustomClaims represents JWT claims structure compatible with serverless auth
type CustomClaims struct {
	UserID   string                 `json:"user_id"`
	UserType string                 `json:"user_type"`
	Custom   map[string]interface{} `json:"custom,omitempty"`
	// Standard JWT fields
	ExpiresAt int64 `json:"exp,omitempty"`
	IssuedAt  int64 `json:"iat,omitempty"`
	NotBefore int64 `json:"nbf,omitempty"`
}

// NewServerlessAuthGateway creates a new serverless authentication gateway
// Similar to NewJWTService from tc-operation-service
func NewServerlessAuthGateway(lambdaAuthURL, serviceAuthURL string) *ServerlessAuthGateway {
	return &ServerlessAuthGateway{
		lambdaAuthURL:  lambdaAuthURL,
		serviceAuthURL: serviceAuthURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// ValidateToken validates JWT token via AWS Lambda ServiceAuth function
// Maintains same interface as JWTService.ValidateToken from tc-operation-service
func (s *ServerlessAuthGateway) ValidateToken(tokenString string) (*CustomClaims, error) {
	if tokenString == "" {
		return nil, fmt.Errorf("token is required")
	}

	// Prepare request payload
	requestPayload := TokenRequest{
		Token: tokenString,
	}

	jsonData, err := json.Marshal(requestPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	// Call Lambda ServiceAuth function
	req, err := http.NewRequest("POST", s.serviceAuthURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "GoLunch-Operation-Service/1.0")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call serverless auth: %v", err)
	}
	defer resp.Body.Close()

	var tokenResponse TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	// Handle different response statuses
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("serverless auth error: %s", tokenResponse.Error)
	}

	if !tokenResponse.Valid {
		return nil, fmt.Errorf("invalid token: %s", tokenResponse.Error)
	}

	if tokenResponse.Claims == nil {
		return nil, fmt.Errorf("no claims returned from serverless auth")
	}

	return tokenResponse.Claims, nil
}

// ValidateServiceToken validates API key for service-to-service communication
// New method specific to microservices architecture
func (s *ServerlessAuthGateway) ValidateServiceToken(apiKey, serviceName string) (bool, error) {
	if apiKey == "" || serviceName == "" {
		return false, fmt.Errorf("api key and service name are required")
	}

	// Prepare service validation request
	serviceRequest := map[string]string{
		"apiKey":      apiKey,
		"serviceName": serviceName,
	}

	jsonData, err := json.Marshal(serviceRequest)
	if err != nil {
		return false, fmt.Errorf("failed to marshal service request: %v", err)
	}

	// Call Lambda ServiceAuth for service validation
	req, err := http.NewRequest("POST", s.serviceAuthURL+"/validate-service", bytes.NewBuffer(jsonData))
	if err != nil {
		return false, fmt.Errorf("failed to create service request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "GoLunch-Operation-Service/1.0")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to call service auth: %v", err)
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK, nil
}
