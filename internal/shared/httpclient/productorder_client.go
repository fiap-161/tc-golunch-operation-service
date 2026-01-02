package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fiap-161/tc-golunch-operation-service/internal/order/entity"
)

type ProductOrderClient struct {
	baseURL string
	client  *http.Client
}

func NewProductOrderClient(baseURL string) *ProductOrderClient {
	return &ProductOrderClient{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

func (c *ProductOrderClient) CreateBulk(ctx context.Context, orderID string, products []entity.OrderProductInfo) error {
	url := fmt.Sprintf("%s/orders/%s/products", c.baseURL, orderID)

	payload := map[string]interface{}{
		"order_id": orderID,
		"products": products,
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
		return fmt.Errorf("failed to create product orders: status %d", resp.StatusCode)
	}

	return nil
}
