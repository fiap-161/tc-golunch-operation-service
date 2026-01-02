package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fiap-161/tc-golunch-operation-service/internal/order/entity"
)

type ProductClient struct {
	baseURL string
	client  *http.Client
}

func NewProductClient(baseURL string) *ProductClient {
	return &ProductClient{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

func (c *ProductClient) FindByIDs(ctx context.Context, productIDs []string) ([]entity.Product, error) {
	url := fmt.Sprintf("%s/admin/products/by-ids", c.baseURL)

	payload := map[string][]string{
		"product_ids": productIDs,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var products []entity.Product
	if err := json.NewDecoder(resp.Body).Decode(&products); err != nil {
		return nil, err
	}

	return products, nil
}
