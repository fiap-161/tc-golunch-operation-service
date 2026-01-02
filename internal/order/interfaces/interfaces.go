package interfaces

import (
	"context"

	"github.com/fiap-161/tc-golunch-operation-service/internal/order/entity"
)

type ProductService interface {
	FindByIDs(ctx context.Context, productIDs []string) ([]entity.Product, error)
}

type ProductOrderService interface {
	CreateBulk(ctx context.Context, orderID string, products []entity.OrderProductInfo) error
}

type PaymentService interface {
	CreateByOrderID(ctx context.Context, orderID string) error
}
