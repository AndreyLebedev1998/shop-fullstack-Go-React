package grpc_package

import (
	"database/sql"

	product "github.com/AndreyLebedev1998/shop-gRPC-product"
	"github.com/redis/go-redis/v9"
)

type Server struct {
	product.UnimplementedProductsServiceServer
	DB  *sql.DB
	RDB *redis.Client
}
