package datasource

import (
	"context"

	"github.com/fiap-161/tc-golunch-operation-service/internal/order/dto"
)

type DataSource interface {
	Create(ctx context.Context, order dto.OrderDAO) (dto.OrderDAO, error)
	GetAll(ctx context.Context) ([]dto.OrderDAO, error)
	FindByID(ctx context.Context, id string) (dto.OrderDAO, error)
	GetPanel(ctx context.Context) ([]dto.OrderDAO, error)
	Update(ctx context.Context, order dto.OrderDAO) (dto.OrderDAO, error)
}
