package constants

type OrderStatus string

type OrderStatusPaid string

const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusConfirmed  OrderStatus = "confirmed"
	OrderStatusProcessing OrderStatus = "processing"
	OrderStatusShipped    OrderStatus = "shipped"
	OrderStatusDelivered  OrderStatus = "delivered"
	OrderStatusCompleted  OrderStatus = "completed"
	OrderStatusCancelled  OrderStatus = "cancelled"
	OrderStatusFailed     OrderStatus = "failed"
)

const (
	OrderStatusPaidCompleted    OrderStatusPaid = "paid"
	OrderStatusPaidNotCompleted OrderStatusPaid = "not paid"
)
