package datasource

import (
	"context"

	"github.com/fiap-161/tc-golunch-operation-service/internal/order/dto"
	"gorm.io/gorm"
)

// DB interface defines the database operations needed
type DB interface {
	Create(value any) *gorm.DB
	Where(query any, args ...any) *gorm.DB
	First(dest any, conds ...any) *gorm.DB
	Find(dest any, conds ...any) *gorm.DB
	Delete(value any, conds ...any) *gorm.DB
	Model(value any) *gorm.DB
	Updates(values any) *gorm.DB
	Save(value any) *gorm.DB
	Order(value any) *gorm.DB
}

// GormDataSource implements DataSource interface using GORM
type GormDataSource struct {
	db DB
}

// New creates a new GormDataSource instance
func New(db DB) DataSource {
	return &GormDataSource{
		db: db,
	}
}

func (g *GormDataSource) Create(ctx context.Context, order dto.OrderDAO) (dto.OrderDAO, error) {
	tx := g.db.Create(&order)
	if tx.Error != nil {
		return dto.OrderDAO{}, tx.Error
	}

	return order, nil
}

func (g *GormDataSource) GetAll(ctx context.Context) ([]dto.OrderDAO, error) {
	var orders []dto.OrderDAO

	if err := g.db.Find(&orders).Error; err != nil {
		return nil, err
	}

	return orders, nil
}

func (g *GormDataSource) FindByID(ctx context.Context, id string) (dto.OrderDAO, error) {
	var order dto.OrderDAO

	tx := g.db.First(&order, "id = ?", id)
	if tx.Error != nil {
		return dto.OrderDAO{}, tx.Error
	}

	return order, nil
}

func (g *GormDataSource) GetPanel(ctx context.Context) ([]dto.OrderDAO, error) {
	var orders []dto.OrderDAO

	if err := g.db.
		Where("status != ?", "completed").
		Order(`
			CASE 
				WHEN status = 'ready' THEN 1
				WHEN status = 'in_preparation' THEN 2
				WHEN status = 'received' THEN 3
				ELSE 4
			END,
			created_at ASC
		`).
		Find(&orders).Error; err != nil {
		return nil, err
	}

	return orders, nil
}

func (g *GormDataSource) Update(ctx context.Context, order dto.OrderDAO) (dto.OrderDAO, error) {
	tx := g.db.Save(&order)
	if tx.Error != nil {
		return dto.OrderDAO{}, tx.Error
	}

	return order, nil
}
