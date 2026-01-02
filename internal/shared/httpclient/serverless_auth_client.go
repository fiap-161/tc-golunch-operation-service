package httpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// ServerlessAuthClient handles communication with AWS Lambda Auth functions
type ServerlessAuthClient struct {
	BaseURL    string
	Client     *http.Client
	APIKey     string
}

// AuthRequest represents the request payload for token validation
type AuthRequest struct {
	Token       string `json:"token"`
	ServiceName string `json:"service_name,omitempty"`
}

// AuthResponse represents the response from Lambda auth functions
type AuthResponse struct {
	Valid   bool                   `json:"valid"`
	Claims  map[string]interface{} `json:"claims,omitempty"`
	UserID  string                 `json:"user_id,omitempty"`
	Role    string                 `json:"role,omitempty"`
	Error   string                 `json:"error,omitempty"`
}

// ServiceAuthRequest represents service-to-service authentication
type ServiceAuthRequest struct {
	FromService string `json:"from_service"`
	ToService   string `json:"to_service"`
	APIKey      string `json:"api_key"`
}

// NewServerlessAuthClient creates a new client for serverless authentication
func NewServerlessAuthClient() *ServerlessAuthClient {
	baseURL := os.Getenv("LAMBDA_AUTH_URL")
	if baseURL == "" {
		baseURL = "https://api-gateway-url.amazonaws.com" // Default fallback
	}

	return &ServerlessAuthClient{
		BaseURL: baseURL,
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
		APIKey: os.Getenv("SERVICE_AUTH_API_KEY"),
	}
}

// ValidateAdminToken validates admin JWT token via Lambda
func (c *ServerlessAuthClient) ValidateAdminToken(token string) (*AuthResponse, error) {
	request := AuthRequest{
		Token: token,
	}

	return c.makeAuthRequest("/admin/validate", request)
}

// ValidateServiceAuth validates service-to-service authentication
func (c *ServerlessAuthClient) ValidateServiceAuth(fromService, toService string) (*AuthResponse, error) {
	request := ServiceAuthRequest{
		FromService: fromService,
		ToService:   toService,
		APIKey:      c.APIKey,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal service auth request: %w", err)
	}

	req, err := http.NewRequest("POST", c.BaseURL+"/service/auth", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create service auth request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", c.APIKey)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make service auth request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read service auth response: %w", err)
	}

	var authResp AuthResponse
	if err := json.Unmarshal(body, &authResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal service auth response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return &authResp, fmt.Errorf("service auth failed: %s", authResp.Error)
	}

	return &authResp, nil
}

// ValidateToken validates any JWT token via Lambda
func (c *ServerlessAuthClient) ValidateToken(token string) (*AuthResponse, error) {
	request := AuthRequest{
		Token: token,
	}

	return c.makeAuthRequest("/auth/validate", request)
}

// makeAuthRequest is a helper method for making authentication requests
func (c *ServerlessAuthClient) makeAuthRequest(endpoint string, request AuthRequest) (*AuthResponse, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal auth request: %w", err)
	}

	req, err := http.NewRequest("POST", c.BaseURL+endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create auth request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.APIKey != "" {
		req.Header.Set("X-API-Key", c.APIKey)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make auth request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read auth response: %w", err)
	}

	var authResp AuthResponse
	if err := json.Unmarshal(body, &authResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal auth response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return &authResp, fmt.Errorf("auth validation failed: %s", authResp.Error)
	}

	return &authResp, nil
}