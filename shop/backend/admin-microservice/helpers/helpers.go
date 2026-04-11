package helpers

import (
	"admin-microservice/constants"
	"slices"
)

func GetAllStatusesOrders() []constants.OrderStatus {
	return []constants.OrderStatus{
		constants.OrderStatusPending,
		constants.OrderStatusConfirmed,
		constants.OrderStatusProcessing,
		constants.OrderStatusShipped,
		constants.OrderStatusDelivered,
		constants.OrderStatusCompleted,
		constants.OrderStatusCancelled,
	}
}

func GetAllStatusesOrdersPaid() []constants.OrderStatusPaid {
	return []constants.OrderStatusPaid{
		constants.OrderStatusPaidCompleted,
		constants.OrderStatusPaidNotCompleted,
	}
}

func IsValidStatus(s constants.OrderStatus) bool {
	return slices.Contains(GetAllStatusesOrders(), s)
}

func IsValidStatusPaid(s constants.OrderStatusPaid) bool {
	return slices.Contains(GetAllStatusesOrdersPaid(), s)
}
