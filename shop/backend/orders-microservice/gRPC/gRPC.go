package grpc_package

import (
	"database/sql"

	"github.com/redis/go-redis/v9"

	product "github.com/AndreyLebedev1998/shop-gRPC-product"
)

type Server struct {
	product.UnimplementedProductsServiceServer
	DB  *sql.DB
	RDB *redis.Client
}
