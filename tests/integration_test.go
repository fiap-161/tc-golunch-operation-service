package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockOrderClient simula o Order Service para testes isolados
type MockOrderClient struct {
	mock.Mock
}

func (m *MockOrderClient) GetOrder(orderID string) (*OrderResponse, error) {
	args := m.Called(orderID)
	return args.Get(0).(*OrderResponse), args.Error(1)
}

func (m *MockOrderClient) UpdateOrderStatus(orderID, status string) error {
	args := m.Called(orderID, status)
	return args.Error(0)
}

// MockPaymentClient simula o Payment Service para testes isolados
type MockPaymentClient struct {
	mock.Mock
}

func (m *MockPaymentClient) GetPayment(paymentID string) (*PaymentResponse, error) {
	args := m.Called(paymentID)
	return args.Get(0).(*PaymentResponse), args.Error(1)
}

// OrderResponse representa resposta do Order Service
type OrderResponse struct {
	ID          string      `json:"id"`
	CustomerID  string      `json:"customer_id"`
	TotalAmount float64     `json:"total_amount"`
	Status      string      `json:"status"`
	Items       []OrderItem `json:"items"`
	CreatedAt   time.Time   `json:"created_at"`
}

// OrderItem representa item do pedido
type OrderItem struct {
	ProductID   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
	TotalPrice  float64 `json:"total_price"`
}

// PaymentResponse representa resposta do Payment Service
type PaymentResponse struct {
	PaymentID string    `json:"payment_id"`
	OrderID   string    `json:"order_id"`
	Amount    float64   `json:"amount"`
	Status    string    `json:"status"`
	PaidAt    time.Time `json:"paid_at"`
}

// TestProductionOrderCreationWithOrderValidation testa criação de ordem de produção com validação mockada
func TestProductionOrderCreationWithOrderValidation(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Mocks
	mockOrderClient := new(MockOrderClient)

	// Mock response do Order Service
	expectedOrder := &OrderResponse{
		ID:          "order_123",
		CustomerID:  "customer_123",
		TotalAmount: 99.90,
		Status:      "paid",
		Items: []OrderItem{
			{ProductID: "prod_1", ProductName: "Hamburger Clássico", Quantity: 2, UnitPrice: 29.90, TotalPrice: 59.80},
			{ProductID: "prod_2", ProductName: "Batata Frita", Quantity: 1, UnitPrice: 15.90, TotalPrice: 15.90},
			{ProductID: "prod_3", ProductName: "Refrigerante", Quantity: 2, UnitPrice: 12.10, TotalPrice: 24.20},
		},
		CreatedAt: time.Now().Add(-30 * time.Minute),
	}

	// Configurar expectativas
	mockOrderClient.On("GetOrder", "order_123").Return(expectedOrder, nil)

	// Simular dados locais de produção
	productionOrders := make(map[string]interface{})

	// Rota para receber notificação de novo pedido pago
	router.POST("/production/orders", func(c *gin.Context) {
		var request struct {
			OrderID   string `json:"order_id"`
			PaymentID string `json:"payment_id"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		// Validar pedido com Order Service (mockado)
		orderResponse, err := mockOrderClient.GetOrder(request.OrderID)
		if err != nil {
			c.JSON(404, gin.H{"error": "Order not found"})
			return
		}

		// Verificar se pedido está pago
		if orderResponse.Status != "paid" {
			c.JSON(400, gin.H{"error": "Order not paid yet"})
			return
		}

		// Criar ordem de produção
		productionOrderID := "prod_order_" + time.Now().Format("20060102150405")

		productionOrder := map[string]interface{}{
			"id":             productionOrderID,
			"order_id":       request.OrderID,
			"customer_id":    orderResponse.CustomerID,
			"status":         "received",
			"priority":       1,
			"estimated_time": calculateEstimatedTime(orderResponse.Items),
			"items":          orderResponse.Items,
			"created_at":     time.Now(),
			"updated_at":     time.Now(),
		}

		productionOrders[productionOrderID] = productionOrder

		c.JSON(201, gin.H{
			"production_order_id": productionOrderID,
			"order_id":            request.OrderID,
			"status":              "received",
			"estimated_time":      productionOrder["estimated_time"],
			"items_count":         len(orderResponse.Items),
		})
	})

	// Teste
	requestData := map[string]interface{}{
		"order_id":   "order_123",
		"payment_id": "payment_123",
	}

	jsonData, _ := json.Marshal(requestData)
	req, _ := http.NewRequest("POST", "/production/orders", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, 201, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.NotEmpty(t, response["production_order_id"])
	assert.Equal(t, "order_123", response["order_id"])
	assert.Equal(t, "received", response["status"])
	assert.NotEmpty(t, response["estimated_time"])
	assert.Equal(t, float64(3), response["items_count"])

	// Verificar mock
	mockOrderClient.AssertExpectations(t)
}

// calculateEstimatedTime calcula tempo estimado baseado nos itens
func calculateEstimatedTime(items []OrderItem) int {
	totalTime := 0
	for _, item := range items {
		// Tempo base por item + tempo por quantidade
		itemTime := 5 + (item.Quantity * 2) // 5 min base + 2 min por unidade
		totalTime += itemTime
	}
	return totalTime
}

// TestProductionStatusUpdate testa atualização de status com notificação mockada
func TestProductionStatusUpdate(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Mocks
	mockOrderClient := new(MockOrderClient)

	// Configurar expectativa
	mockOrderClient.On("UpdateOrderStatus", "order_123", "in_production").Return(nil)

	// Simular dados locais
	productionOrders := map[string]interface{}{
		"prod_order_123": map[string]interface{}{
			"id":         "prod_order_123",
			"order_id":   "order_123",
			"status":     "received",
			"created_at": time.Now().Add(-10 * time.Minute),
		},
	}

	// Rota para atualizar status
	router.PUT("/production/orders/:id/status", func(c *gin.Context) {
		productionOrderID := c.Param("id")

		var request struct {
			Status string `json:"status"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		// Verificar se ordem existe
		productionOrder, exists := productionOrders[productionOrderID]
		if !exists {
			c.JSON(404, gin.H{"error": "Production order not found"})
			return
		}

		prodOrderMap := productionOrder.(map[string]interface{})
		orderID := prodOrderMap["order_id"].(string)

		// Atualizar status local
		prodOrderMap["status"] = request.Status
		prodOrderMap["updated_at"] = time.Now()

		// Mapear status da produção para status do pedido
		var orderStatus string
		switch request.Status {
		case "in_preparation":
			orderStatus = "in_production"
		case "ready":
			orderStatus = "ready"
		case "delivered":
			orderStatus = "completed"
		default:
			orderStatus = request.Status
		}

		// Notificar Order Service (mockado)
		if orderStatus != request.Status {
			err := mockOrderClient.UpdateOrderStatus(orderID, orderStatus)
			if err != nil {
				c.JSON(500, gin.H{"error": "Failed to update order status"})
				return
			}
		}

		c.JSON(200, gin.H{
			"production_order_id": productionOrderID,
			"order_id":            orderID,
			"status":              request.Status,
			"order_status":        orderStatus,
			"updated_at":          prodOrderMap["updated_at"],
		})
	})

	// Teste
	statusData := map[string]interface{}{
		"status": "in_preparation",
	}

	jsonData, _ := json.Marshal(statusData)
	req, _ := http.NewRequest("PUT", "/production/orders/prod_order_123/status", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, 200, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "prod_order_123", response["production_order_id"])
	assert.Equal(t, "order_123", response["order_id"])
	assert.Equal(t, "in_preparation", response["status"])
	assert.Equal(t, "in_production", response["order_status"])

	// Verificar mock
	mockOrderClient.AssertExpectations(t)
}

// TestProductionPanelWithoutExternalDependencies testa painel de produção sem dependências externas
func TestProductionPanelWithoutExternalDependencies(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Mock de dados locais
	productionOrders := []map[string]interface{}{
		{
			"id":             "prod_order_1",
			"order_id":       "order_123",
			"customer_id":    "customer_123",
			"status":         "received",
			"priority":       1,
			"estimated_time": 15,
			"created_at":     time.Now().Add(-5 * time.Minute),
			"items_count":    3,
		},
		{
			"id":             "prod_order_2",
			"order_id":       "order_456",
			"customer_id":    "customer_456",
			"status":         "in_preparation",
			"priority":       2,
			"estimated_time": 10,
			"created_at":     time.Now().Add(-15 * time.Minute),
			"started_at":     time.Now().Add(-10 * time.Minute),
			"items_count":    2,
		},
		{
			"id":             "prod_order_3",
			"order_id":       "order_789",
			"customer_id":    "customer_789",
			"status":         "ready",
			"priority":       1,
			"estimated_time": 8,
			"created_at":     time.Now().Add(-25 * time.Minute),
			"completed_at":   time.Now().Add(-5 * time.Minute),
			"items_count":    1,
		},
	}

	// Rota do painel (requer autenticação - simulada)
	router.GET("/admin/orders/panel", func(c *gin.Context) {
		status := c.Query("status")

		var filteredOrders []map[string]interface{}

		for _, order := range productionOrders {
			if status == "" || order["status"] == status {
				filteredOrders = append(filteredOrders, order)
			}
		}

		// Estatísticas
		stats := map[string]int{
			"total":          len(filteredOrders),
			"received":       0,
			"in_preparation": 0,
			"ready":          0,
			"completed":      0,
		}

		for _, order := range filteredOrders {
			orderStatus := order["status"].(string)
			stats[orderStatus]++
		}

		c.JSON(200, gin.H{
			"orders": filteredOrders,
			"stats":  stats,
		})
	})

	// Teste sem filtro
	req, _ := http.NewRequest("GET", "/admin/orders/panel", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	ordersResponse := response["orders"].([]interface{})
	assert.Len(t, ordersResponse, 3)

	stats := response["stats"].(map[string]interface{})
	assert.Equal(t, float64(3), stats["total"])
	assert.Equal(t, float64(1), stats["received"])
	assert.Equal(t, float64(1), stats["in_preparation"])
	assert.Equal(t, float64(1), stats["ready"])

	// Teste com filtro
	req, _ = http.NewRequest("GET", "/admin/orders/panel?status=in_preparation", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	ordersResponse = response["orders"].([]interface{})
	assert.Len(t, ordersResponse, 1)

	firstOrder := ordersResponse[0].(map[string]interface{})
	assert.Equal(t, "in_preparation", firstOrder["status"])
}

// TestAdminAuthentication testa autenticação de admin sem dependências externas
func TestAdminAuthentication(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Mock de dados de admin
	admins := map[string]map[string]interface{}{
		"testchef": {
			"id":            "admin_123",
			"name":          "Test Chef",
			"login":         "testchef",
			"password_hash": "hashed_password_123", // Simular hash
			"role":          "chef",
			"active":        true,
		},
	}

	// Rota de login
	router.POST("/admin/login", func(c *gin.Context) {
		var loginRequest struct {
			Login    string `json:"login"`
			Password string `json:"password"`
		}

		if err := c.ShouldBindJSON(&loginRequest); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		// Verificar credenciais (simulado)
		admin, exists := admins[loginRequest.Login]
		if !exists {
			c.JSON(401, gin.H{"error": "Invalid credentials"})
			return
		}

		// Simular verificação de senha (normalmente seria hash)
		if loginRequest.Password != "test123" {
			c.JSON(401, gin.H{"error": "Invalid credentials"})
			return
		}

		if !admin["active"].(bool) {
			c.JSON(401, gin.H{"error": "Admin account disabled"})
			return
		}

		// Simular geração de JWT token
		token := "jwt_token_" + time.Now().Format("20060102150405")

		c.JSON(200, gin.H{
			"token":      token,
			"admin_id":   admin["id"],
			"name":       admin["name"],
			"role":       admin["role"],
			"expires_in": 3600, // 1 hora
		})
	})

	// Teste de login válido
	loginData := map[string]interface{}{
		"login":    "testchef",
		"password": "test123",
	}

	jsonData, _ := json.Marshal(loginData)
	req, _ := http.NewRequest("POST", "/admin/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.NotEmpty(t, response["token"])
	assert.Equal(t, "admin_123", response["admin_id"])
	assert.Equal(t, "Test Chef", response["name"])
	assert.Equal(t, "chef", response["role"])
	assert.Equal(t, float64(3600), response["expires_in"])

	// Teste de login inválido
	loginData["password"] = "wrong_password"
	jsonData, _ = json.Marshal(loginData)
	req, _ = http.NewRequest("POST", "/admin/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
}

// TestProductionMetrics testa métricas de produção sem dependências externas
func TestProductionMetrics(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Mock de dados históricos
	completedOrders := []map[string]interface{}{
		{
			"order_id":       "order_1",
			"estimated_time": 15,
			"actual_time":    12,
			"completed_at":   time.Now().Add(-1 * time.Hour),
		},
		{
			"order_id":       "order_2",
			"estimated_time": 10,
			"actual_time":    15,
			"completed_at":   time.Now().Add(-2 * time.Hour),
		},
		{
			"order_id":       "order_3",
			"estimated_time": 8,
			"actual_time":    8,
			"completed_at":   time.Now().Add(-3 * time.Hour),
		},
	}

	router.GET("/admin/metrics", func(c *gin.Context) {
		totalOrders := len(completedOrders)
		totalEstimated := 0
		totalActual := 0

		for _, order := range completedOrders {
			totalEstimated += int(order["estimated_time"].(int))
			totalActual += int(order["actual_time"].(int))
		}

		avgEstimated := float64(totalEstimated) / float64(totalOrders)
		avgActual := float64(totalActual) / float64(totalOrders)
		efficiency := (avgEstimated / avgActual) * 100

		c.JSON(200, gin.H{
			"total_orders":          totalOrders,
			"avg_estimated_time":    avgEstimated,
			"avg_actual_time":       avgActual,
			"efficiency_percentage": efficiency,
			"orders":                completedOrders,
		})
	})

	// Teste
	req, _ := http.NewRequest("GET", "/admin/metrics", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, float64(3), response["total_orders"])
	assert.Equal(t, 11.0, response["avg_estimated_time"])            // (15+10+8)/3
	assert.Equal(t, 11.666666666666666, response["avg_actual_time"]) // (12+15+8)/3
	assert.NotEmpty(t, response["efficiency_percentage"])

	ordersResponse := response["orders"].([]interface{})
	assert.Len(t, ordersResponse, 3)
}
