package datasource

import (
	"context"

	"github.com/fiap-161/tc-golunch-operation-service/internal/admin/dto"
	"gorm.io/gorm"
)

type DB interface {
	Create(value any) *gorm.DB
	Where(query any, args ...any) *gorm.DB
	First(dest any, conds ...any) *gorm.DB
}

type GormDataSource struct {
	db DB
}

func New(db DB) *GormDataSource {
	return &GormDataSource{
		db: db,
	}
}

func (r *GormDataSource) Create(_ context.Context, admin dto.AdminDAO) error {
	tx := r.db.Create(&admin)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (r *GormDataSource) FindByEmail(_ context.Context, email string) (dto.AdminDAO, error) {
	var admin dto.AdminDAO

	tx := r.db.Where("email = ?", email).First(&admin)

	if tx.Error != nil {
		return dto.AdminDAO{}, tx.Error
	}

	return admin, nil
}
