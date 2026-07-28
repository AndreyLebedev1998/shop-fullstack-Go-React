package grpc_package

import (
	"context"
	"database/sql"
	"fmt"
	grpcFunctionsQuery "orders-microservice/gRPC-functions-query"
	"orders-microservice/models"

	"github.com/redis/go-redis/v9"

	order "github.com/AndreyLebedev1998/shop-gRPC-orders"
	product "github.com/AndreyLebedev1998/shop-gRPC-product"
)

type Server struct {
	product.UnimplementedProductsServiceServer
	order.UnimplementedOrderServiceServer
	DB  *sql.DB
	RDB *redis.Client
}

func (s *Server) ChangeOrderStatus(ctx context.Context, orderSatus *order.OrdersStatus) (*order.OrderStatusResponse, error) {
	o := models.OrderStatus{
		Id:     orderSatus.Id,
		Status: orderSatus.Status,
	}
	message, err := grpcFunctionsQuery.ChangeOrderStatus(s.DB, o)
	if err != nil {
		return nil, err
	}

	return &order.OrderStatusResponse{
		Response: message.Response,
		Status:   message.Status,
		Id:       message.Id,
	}, nil
}

func (s *Server) ChangeOrderStatusPaid(ctx context.Context, orderSatus *order.OrderStatusPaid) (*order.OrderStatusResponse, error) {
	o := models.OrderStatusPaid{
		Id:         orderSatus.Id,
		StatusPaid: orderSatus.StatusPaid,
	}
	message, err := grpcFunctionsQuery.ChangeOrderStatusPaid(s.DB, o)
	if err != nil {
		return nil, err
	}

	return &order.OrderStatusResponse{
		Response: message.Response,
		Status:   message.Status,
		Id:       message.Id,
	}, nil
}

func (s *Server) RemoveOrder(ctx context.Context, orderId *order.OrderId) (*order.OrderMessage, error) {
	var statusCancelled = "cancelled"
	query := "UPDATE orders SET status = $1 WHERE id = $2"
	queryOrder := "SELECT status FROM orders WHERE id = $1"

	var status string

	errGet := s.DB.QueryRow(queryOrder, orderId.Id).Scan(&status)
	if errGet != nil {
		if errGet == sql.ErrNoRows {
			return nil, fmt.Errorf("order not found")
		}
		return nil, fmt.Errorf("server error: %v", errGet)
	}

	if status == "cancelled" {
		return nil, fmt.Errorf("the order has already been canceled")
	}

	rows, err := s.DB.Exec(query, statusCancelled, orderId.Id)
	if err != nil {
		return nil, fmt.Errorf("error update order status")
	}

	rowsUpdated, _ := rows.RowsAffected()

	if rowsUpdated != 1 {
		return nil, fmt.Errorf("order id is not defined")
	}

	return &order.OrderMessage{
		Message: "order is cancelled",
		Id:      orderId.Id,
		Status:  true,
	}, nil
}

func (s *Server) GetOrderItems(ctx context.Context, orderId *order.OrderId) (*order.OrderItems, error) {
	queryOrderItems := "SELECT id, product_id, quantity FROM order_items WHERE order_id = $1"

	var orderItems []models.OrderItems
	rows, err := s.DB.Query(queryOrderItems, orderId.Id)
	if err != nil {
		return nil, fmt.Errorf("Server error")
	}

	for rows.Next() {
		var orderItem models.OrderItems

		err := rows.Scan(&orderItem.Id, &orderItem.ProductId, &orderItem.Quantity)
		if err != nil {
			return nil, fmt.Errorf("error scan order_items")
		}

		orderItems = append(orderItems, orderItem)
	}

	var orderItemsGRPC []*order.OrderItem

	for _, o := range orderItems {
		orderItemsGRPC = append(orderItemsGRPC, &order.OrderItem{
			Id:        o.Id,
			ProductId: o.ProductId,
			Quantity:  o.Quantity,
		})
	}

	return &order.OrderItems{
		Items: orderItemsGRPC,
	}, nil
}

func (s *Server) GetProductsIdsFromOrdersForUsers(ctx context.Context, userId *order.UserId) (*order.ProductIds, error) {
	query := "SELECT product_id FROM order_items JOIN orders ON order_items.order_id = orders.id WHERE user_id = $1 LIMIT 25"

	var productIds []*order.ProductId = make([]*order.ProductId, 0)

	rows, err := s.DB.QueryContext(ctx, query, userId.UserId)

	if err == sql.ErrNoRows {
		return &order.ProductIds{
			ProductsIds: productIds,
		}, nil
	}

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var productId int64

		err := rows.Scan(&productId)
		if err != nil {
			return nil, err
		}

		productIds = append(productIds, &order.ProductId{
			Id: productId,
		})
	}

	return &order.ProductIds{
		ProductsIds: productIds,
	}, nil
}
