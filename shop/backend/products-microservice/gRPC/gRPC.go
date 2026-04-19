package grpc_package

import (
	"context"
	"database/sql"
	"errors"
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
	if len(products.NewItems) == 0 || len(products.OldItems) == 0 {
		return nil, fmt.Errorf("products is empty")
	}

	newItems := make([]models.UpdateProductForGRPC, 0, len(products.NewItems))
	oldItems := make([]models.UpdateProductForGRPC, 0, len(products.OldItems))

	for _, p := range products.NewItems {
		newItems = append(newItems, models.UpdateProductForGRPC{
			ProductId: p.ProductId,
			Quantity:  p.Quantity,
		})
	}

	for _, p := range products.OldItems {
		oldItems = append(oldItems, models.UpdateProductForGRPC{
			ProductId: p.ProductId,
			Quantity:  p.Quantity,
		})
	}

	message, err := grpcFunctionsQuery.UpdateProductsQuantityByIds(s.DB, newItems, oldItems, products.IsCreate)
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

func (s *Server) CreateProduct(ctx context.Context, newProduct *product.NewProduct) (*product.ReturnNewProduct, error) {
	query := `INSERT INTO products (product_name, price, category_id, image_url, availability_of_pieces) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	var id int

	err := s.DB.QueryRow(query, newProduct.ProductName, newProduct.Price, newProduct.CategoryId, newProduct.ImageUrl, newProduct.AvailabilityOfPieces).Scan(&id)
	if err != nil {
		return nil, errors.New("error INSERT INTO products")
	}

	return &product.ReturnNewProduct{
		Id:                   int64(id),
		ProductName:          newProduct.ProductName,
		Price:                newProduct.Price,
		CategoryId:           newProduct.CategoryId,
		ImageUrl:             newProduct.ImageUrl,
		AvailabilityOfPieces: newProduct.AvailabilityOfPieces,
	}, nil
}
