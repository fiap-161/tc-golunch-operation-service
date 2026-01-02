package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type PaymentClient struct {
	baseURL string
	client  *http.Client
}

func NewPaymentClient(baseURL string) *PaymentClient {
	return &PaymentClient{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

func (c *PaymentClient) CreateByOrderID(ctx context.Context, orderID string) error {
	url := fmt.Sprintf("%s/payments", c.baseURL)

	payload := map[string]string{
		"order_id": orderID,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("failed to create payment: status %d", resp.StatusCode)
	}

	return nil
}
