package grpc_package

import (
	"context"
	"database/sql"
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
