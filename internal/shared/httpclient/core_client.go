package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type CoreServiceClient struct {
	baseURL    string
	httpClient *http.Client
}

type OrderResponse struct {
	ID       string `json:"id"`
	Status   string `json:"status"`
	Customer struct {
		Name string `json:"name"`
		CPF  string `json:"cpf"`
	} `json:"customer"`
	Products []struct {
		Name     string  `json:"name"`
		Quantity int     `json:"quantity"`
		Price    float64 `json:"price"`
	} `json:"products"`
	TotalAmount float64   `json:"total_amount"`
	CreatedAt   time.Time `json:"created_at"`
}

type OrderUpdateRequest struct {
	Status string `json:"status"`
}

func NewCoreServiceClient() *CoreServiceClient {
	baseURL := os.Getenv("CORE_SERVICE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8081"
	}

	return &CoreServiceClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *CoreServiceClient) GetOrder(ctx context.Context, orderID string) (*OrderResponse, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/order/"+orderID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("core service returned status %d", resp.StatusCode)
	}

	var orderResp OrderResponse
	if err := json.NewDecoder(resp.Body).Decode(&orderResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &orderResp, nil
}

func (c *CoreServiceClient) GetAllOrders(ctx context.Context) ([]OrderResponse, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/order", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("core service returned status %d", resp.StatusCode)
	}

	var orders []OrderResponse
	if err := json.NewDecoder(resp.Body).Decode(&orders); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return orders, nil
}

func (c *CoreServiceClient) UpdateOrderStatus(ctx context.Context, orderID, status string) error {
	payload := OrderUpdateRequest{Status: status}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal update request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "PUT", c.baseURL+"/order/"+orderID, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("core service returned status %d", resp.StatusCode)
	}

	return nil
}
