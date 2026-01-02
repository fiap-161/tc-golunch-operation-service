package gateway

import (
	"context"

	"github.com/fiap-161/tc-golunch-operation-service/internal/order/dto"
	"github.com/fiap-161/tc-golunch-operation-service/internal/order/entity"
	"github.com/fiap-161/tc-golunch-operation-service/internal/order/external/datasource"
	apperror "github.com/fiap-161/tc-golunch-operation-service/internal/shared/errors"
)

type Gateway struct {
	Datasource datasource.DataSource
}

func Build(datasource datasource.DataSource) *Gateway {
	return &Gateway{
		Datasource: datasource,
	}
}

func (g *Gateway) Create(ctx context.Context, order entity.Order) (entity.Order, error) {
	orderDAO := dto.ToOrderDAO(order)
	created, err := g.Datasource.Create(ctx, orderDAO)
	if err != nil {
		return entity.Order{}, &apperror.InternalError{Msg: err.Error()}
	}
	return dto.FromOrderDAO(created), nil
}

func (g *Gateway) GetAll(ctx context.Context) ([]entity.Order, error) {
	ordersDAO, err := g.Datasource.GetAll(ctx)
	if err != nil {
		return nil, &apperror.InternalError{Msg: err.Error()}
	}
	return dto.EntityListFromDAOList(ordersDAO), nil
}

func (g *Gateway) GetPanel(ctx context.Context) ([]entity.Order, error) {
	ordersDAO, err := g.Datasource.GetPanel(ctx)
	if err != nil {
		return nil, &apperror.InternalError{Msg: err.Error()}
	}
	return dto.EntityListFromDAOList(ordersDAO), nil
}

func (g *Gateway) FindByID(ctx context.Context, id string) (entity.Order, error) {
	orderDAO, err := g.Datasource.FindByID(ctx, id)
	if err != nil {
		return entity.Order{}, err
	}
	return dto.FromOrderDAO(orderDAO), nil
}

func (g *Gateway) Update(ctx context.Context, order entity.Order) (entity.Order, error) {
	orderDAO := dto.ToOrderDAO(order)
	updated, err := g.Datasource.Update(ctx, orderDAO)
	if err != nil {
		return entity.Order{}, &apperror.InternalError{Msg: err.Error()}
	}
	return dto.FromOrderDAO(updated), nil
}
