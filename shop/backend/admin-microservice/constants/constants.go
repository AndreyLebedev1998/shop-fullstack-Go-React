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
)

const (
	OrderStatusPaidCompleted    OrderStatusPaid = "paid"
	OrderStatusPaidNotCompleted OrderStatusPaid = "not_paid"
)
