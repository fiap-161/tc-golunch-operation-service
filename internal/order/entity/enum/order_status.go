package enum

type OrderStatus string

const (
	OrderStatusAwaitingPayment OrderStatus = "awaiting_payment"
	OrderStatusReceived        OrderStatus = "received"
	OrderStatusInPreparation   OrderStatus = "in_preparation"
	OrderStatusReady           OrderStatus = "ready"
	OrderStatusCompleted       OrderStatus = "completed"
)

var OrderPanelStatus = []string{
	OrderStatusReceived.String(),
	OrderStatusInPreparation.String(),
	OrderStatusReady.String(),
}

var StatusMapper = map[string]OrderStatus{
	OrderStatusAwaitingPayment.String(): OrderStatusAwaitingPayment,
	OrderStatusReceived.String():        OrderStatusReceived,
	OrderStatusInPreparation.String():   OrderStatusInPreparation,
	OrderStatusReady.String():           OrderStatusReady,
	OrderStatusCompleted.String():       OrderStatusCompleted,
}

func (o OrderStatus) String() string {
	return string(o)
}
