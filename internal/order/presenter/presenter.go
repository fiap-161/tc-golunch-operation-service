package presenter

import (
	"github.com/fiap-161/tc-golunch-operation-service/internal/order/dto"
	"github.com/fiap-161/tc-golunch-operation-service/internal/order/entity"
)

type Presenter struct{}

func Build() *Presenter {
	return &Presenter{}
}

func (p *Presenter) FromEntityToDAO(order entity.Order) dto.OrderDAO {
	return dto.ToOrderDAO(order)
}

func (p *Presenter) FromEntityListToDAOList(orders []entity.Order) []dto.OrderDAO {
	var ordersDAO []dto.OrderDAO
	for _, order := range orders {
		ordersDAO = append(ordersDAO, dto.ToOrderDAO(order))
	}
	return ordersDAO
}
