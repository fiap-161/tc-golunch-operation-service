package usecases

import (
	"context"

	"github.com/fiap-161/tc-golunch-operation-service/internal/order/dto"
	"github.com/fiap-161/tc-golunch-operation-service/internal/order/entity"
	"github.com/fiap-161/tc-golunch-operation-service/internal/order/gateway"
	"github.com/fiap-161/tc-golunch-operation-service/internal/order/interfaces"
	apperror "github.com/fiap-161/tc-golunch-operation-service/internal/shared/errors"
)

type UseCases struct {
	orderGateway        *gateway.Gateway
	productService      interfaces.ProductService
	productOrderService interfaces.ProductOrderService
	paymentService      interfaces.PaymentService
}

func Build(
	orderGateway *gateway.Gateway,
	productService interfaces.ProductService,
	productOrderService interfaces.ProductOrderService,
	paymentService interfaces.PaymentService,
) *UseCases {
	return &UseCases{
		orderGateway:        orderGateway,
		productService:      productService,
		productOrderService: productOrderService,
		paymentService:      paymentService,
	}
}

func (u *UseCases) CreateCompleteOrder(ctx context.Context, orderDTO dto.CreateOrderDTO) (string, error) {
	var productIds []string
	for _, item := range orderDTO.Products {
		productIds = append(productIds, item.ProductID)
	}

	products, findErr := u.productService.FindByIDs(ctx, productIds)
	if findErr != nil {
		return "", findErr
	}
	if len(products) != len(orderDTO.Products) {
		return "", &apperror.NotFoundError{
			Msg: "some products not found",
		}
	}

	// Criar pedido
	populatedOrder := generateOrderByProducts(orderDTO, products)
	createdOrder, createErr := u.orderGateway.Create(ctx, populatedOrder.Build())
	if createErr != nil {
		return "", createErr
	}

	// Converter para entity.OrderProductInfo para a interface
	orderProductInfo := make([]entity.OrderProductInfo, len(orderDTO.Products))
	for i, product := range orderDTO.Products {
		orderProductInfo[i] = entity.OrderProductInfo{
			ProductID: product.ProductID,
			Quantity:  product.Quantity,
		}
	}

	createBulkErr := u.productOrderService.CreateBulk(ctx, createdOrder.ID, orderProductInfo)
	if createBulkErr != nil {
		return "", createBulkErr
	}

	paymentErr := u.paymentService.CreateByOrderID(ctx, createdOrder.ID)
	if paymentErr != nil {
		return "", paymentErr
	}

	return "payment-qr-code-placeholder", nil
}

func generateOrderByProducts(orderDTO dto.CreateOrderDTO, products []entity.Product) entity.Order {
	orderProductInfo := make([]entity.OrderProductInfo, len(orderDTO.Products))
	for i, product := range orderDTO.Products {
		orderProductInfo[i] = entity.OrderProductInfo{
			ProductID: product.ProductID,
			Quantity:  product.Quantity,
		}
	}

	return entity.Order{}.FromDTO(orderDTO.CustomerID, orderProductInfo, products)
}
func (u *UseCases) CreateOrder(ctx context.Context, order entity.Order) (entity.Order, error) {
	return u.orderGateway.Create(ctx, order)
}

func (u *UseCases) GetAllOrById(ctx context.Context, id string) ([]entity.Order, error) {
	if id != "" {
		order, err := u.orderGateway.FindByID(ctx, id)
		if err != nil {
			return nil, err
		}
		return []entity.Order{order}, nil
	}
	return u.orderGateway.GetAll(ctx)
}

func (u *UseCases) GetPanel(ctx context.Context) ([]entity.Order, error) {
	return u.orderGateway.GetPanel(ctx)
}

func (u *UseCases) FindByID(ctx context.Context, id string) (entity.Order, error) {
	return u.orderGateway.FindByID(ctx, id)
}

func (u *UseCases) Update(ctx context.Context, order entity.Order) (entity.Order, error) {
	return u.orderGateway.Update(ctx, order)
}
