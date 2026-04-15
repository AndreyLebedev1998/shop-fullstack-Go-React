package grpc_package

import (
	"context"
	"database/sql"
	"fmt"
	grpcFunctionsQuery "products-microservice/gRPC-functions-query"
	"products-microservice/helpers"
	"products-microservice/models"

	product "github.com/AndreyLebedev1998/shop-gRPC-product"

	"github.com/redis/go-redis/v9"
)

type Server struct {
	product.UnimplementedProductsServiceServer
	DB  *sql.DB
	RDB *redis.Client
}

func (s *Server) GetProductsByIds(ctx context.Context, products_ids *product.GetProductsRequest) (*product.GetProductsResponse, error) {
	fmt.Println(products_ids)
	if len(products_ids.ProductIds) == 0 {
		return nil, fmt.Errorf("product_ids is empty")
	}

	products, err := grpcFunctionsQuery.GetProductByIds(s.DB, helpers.ConvertInt64ToInt(products_ids.ProductIds))
	if err != nil {
		return nil, fmt.Errorf("server error")
	}

	return &product.GetProductsResponse{
		Products: helpers.ConvertProductsToProto(products),
	}, nil
}

func (s *Server) UpdateProductQuantityByIds(ctx context.Context, products *product.UpdateProductQuantityRequest) (*product.UpdateProductQuantityResponse, error) {
	if len(products.Items) == 0 {
		return nil, fmt.Errorf("products is empty")
	}

	internalProducts := make([]models.UpdateProductForGRPC, 0, len(products.Items))

	for _, p := range products.Items {
		internalProducts = append(internalProducts, models.UpdateProductForGRPC{
			ProductId: p.ProductId,
			Quantity:  p.Quantity,
		})
	}

	message, err := grpcFunctionsQuery.UpdateProductsQuantityByIds(s.DB, internalProducts)
	if err != nil {
		return &product.UpdateProductQuantityResponse{
			Success: message.Success,
			Message: message.Message,
		}, err
	}

	return &product.UpdateProductQuantityResponse{
		Success: message.Success,
		Message: message.Message,
	}, nil
}
