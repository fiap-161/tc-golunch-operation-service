package services

import (
	"context"

	"github.com/fiap-161/tc-golunch-operation-service/internal/order/entity"
	"github.com/fiap-161/tc-golunch-operation-service/internal/order/usecases"
)

type OrderServiceGateway struct {
	orderUseCase *usecases.UseCases
}

func NewOrderServiceGateway(orderUseCase *usecases.UseCases) *OrderServiceGateway {
	return &OrderServiceGateway{
		orderUseCase: orderUseCase,
	}
}

func (a *OrderServiceGateway) FindByID(ctx context.Context, orderID string) (entity.Order, error) {
	order, err := a.orderUseCase.FindByID(ctx, orderID)
	if err != nil {
		return entity.Order{}, err
	}

	return entity.Order{
		Entity:        order.Entity,
		CustomerID:    order.Entity.ID,
		Status:        order.Status,
		Price:         order.Price,
		PreparingTime: order.PreparingTime,
	}, nil
}

func (a *OrderServiceGateway) Update(ctx context.Context, order entity.Order) (entity.Order, error) {
	currentOrder, err := a.orderUseCase.FindByID(ctx, order.Entity.ID)
	if err != nil {
		return entity.Order{}, err
	}

	currentOrder.Status = order.Status

	updatedOrder, updateErr := a.orderUseCase.Update(ctx, currentOrder)
	if updateErr != nil {
		return entity.Order{}, updateErr
	}

	return entity.Order{
		Entity:        updatedOrder.Entity,
		CustomerID:    updatedOrder.CustomerID,
		Status:        updatedOrder.Status,
		Price:         updatedOrder.Price,
		PreparingTime: updatedOrder.PreparingTime,
	}, nil
}
