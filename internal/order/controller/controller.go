package controller

import (
	"context"

	"github.com/fiap-161/tc-golunch-operation-service/internal/order/dto"
	"github.com/fiap-161/tc-golunch-operation-service/internal/order/presenter"
	"github.com/fiap-161/tc-golunch-operation-service/internal/order/usecases"
)

type Controller struct {
	orderUseCase *usecases.UseCases
}

func Build(orderUseCase *usecases.UseCases) *Controller {
	return &Controller{
		orderUseCase: orderUseCase,
	}
}

func (c *Controller) Create(ctx context.Context, orderDTO dto.CreateOrderDTO) (string, error) {
	return c.orderUseCase.CreateCompleteOrder(ctx, orderDTO)
}

func (c *Controller) GetAll(ctx context.Context, id string) ([]dto.OrderDAO, error) {
	presenter := presenter.Build()

	orders, err := c.orderUseCase.GetAllOrById(ctx, id)
	if err != nil {
		return nil, err
	}

	return presenter.FromEntityListToDAOList(orders), nil
}

func (c *Controller) GetPanel(ctx context.Context) ([]dto.OrderDAO, error) {
	presenter := presenter.Build()

	orders, err := c.orderUseCase.GetPanel(ctx)
	if err != nil {
		return nil, err
	}

	return presenter.FromEntityListToDAOList(orders), nil
}

func (c *Controller) FindByID(ctx context.Context, id string) (dto.OrderDAO, error) {
	presenter := presenter.Build()

	order, err := c.orderUseCase.FindByID(ctx, id)
	if err != nil {
		return dto.OrderDAO{}, err
	}

	return presenter.FromEntityToDAO(order), nil
}

func (c *Controller) Update(ctx context.Context, orderDTO dto.OrderDAO) (dto.OrderDAO, error) {
	presenter := presenter.Build()

	order := dto.FromOrderDAO(orderDTO)
	updated, err := c.orderUseCase.Update(ctx, order)
	if err != nil {
		return dto.OrderDAO{}, err
	}

	return presenter.FromEntityToDAO(updated), nil
}
