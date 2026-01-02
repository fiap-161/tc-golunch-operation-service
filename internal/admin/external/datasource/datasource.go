package datasource

import (
	"context"

	"github.com/fiap-161/tc-golunch-operation-service/internal/admin/dto"
)

type DataSource interface {
	Create(ctx context.Context, admin dto.AdminDAO) error
	FindByEmail(ctx context.Context, email string) (dto.AdminDAO, error)
}
