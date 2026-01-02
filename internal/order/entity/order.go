package entity

import (
	"time"

	"github.com/fiap-161/tc-golunch-operation-service/internal/order/entity/enum"
	"github.com/fiap-161/tc-golunch-operation-service/internal/shared/entity"
	"github.com/google/uuid"
)

type Order struct {
	entity.Entity
	CustomerID    string           `json:"customer_id" gorm:"index"`
	Status        enum.OrderStatus `json:"status" gorm:"type:varchar(20)"`
	Price         float64          `json:"price" gorm:"type:decimal(10,2)"`
	PreparingTime uint             `json:"preparing_time" gorm:"type:integer"`
}

func (o Order) Build() Order {
	return Order{
		Entity: entity.Entity{
			ID:        uuid.NewString(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		CustomerID:    o.CustomerID,
		Status:        o.Status,
		Price:         o.Price,
		PreparingTime: o.PreparingTime,
	}
}

type OrderProductInfo struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

// Product representa os dados básicos de um produto para cálculos de pedido
type Product struct {
	Id            string  `json:"id"`
	Price         float64 `json:"price"`
	PreparingTime uint    `json:"preparing_time"`
}

func (o Order) FromDTO(customerID string, products []OrderProductInfo, allProducts []Product) Order {
	totalPrice, preparingTime := o.getOrderInfoFromProducts(allProducts, products)

	return Order{
		CustomerID:    customerID,
		Price:         totalPrice,
		PreparingTime: preparingTime,
		Status:        enum.OrderStatusAwaitingPayment,
	}
}

func (o Order) getOrderInfoFromProducts(products []Product, orderProducts []OrderProductInfo) (float64, uint) {
	var totalPrice float64
	var preparingTime uint

	for _, item := range orderProducts {
		for _, product := range products {
			if product.Id == item.ProductID {
				totalPrice += product.Price * float64(item.Quantity)
				preparingTime += product.PreparingTime
			}
		}
	}

	return totalPrice, preparingTime
}
