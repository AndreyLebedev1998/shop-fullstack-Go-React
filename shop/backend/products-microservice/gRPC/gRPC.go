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

	message, err := grpcFunctionsQuery.UpdateProductsQuantityByIds(s.DB, newItems, oldItems, products.IsCreate, products.IsDelete)
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
	query := `INSERT INTO products (product_name, price, category_id, image_url, availability_of_pieces, subcategory_id) VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT (product_name)
			  DO UPDATE SET category_id = EXCLUDED.category_id, price = EXCLUDED.price, image_url = EXCLUDED.image_url, availability_of_pieces = EXCLUDED.availability_of_pieces, 
			  subcategory_id = EXCLUDED.subcategory_id RETURNING id`
	var id int

	err := s.DB.QueryRow(query, newProduct.ProductName, newProduct.Price, newProduct.CategoryId, newProduct.ImageUrl, newProduct.AvailabilityOfPieces, newProduct.SubcategoryId).Scan(&id)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	return &product.ReturnNewProduct{
		Id:                   int64(id),
		ProductName:          newProduct.ProductName,
		Price:                newProduct.Price,
		CategoryId:           newProduct.CategoryId,
		ImageUrl:             newProduct.ImageUrl,
		AvailabilityOfPieces: newProduct.AvailabilityOfPieces,
		SubcategoryId:        newProduct.SubcategoryId,
	}, nil
}

func (s *Server) CreateCategory(ctx context.Context, category_name *product.NewCategory) (*product.ReturnNewCategory, error) {
	query := "INSERT INTO categories (category_name) VALUES ($1) RETURNING id"

	var category_id int64

	err := s.DB.QueryRow(query, category_name.CategoryName).Scan(&category_id)
	if err != nil {
		return nil, fmt.Errorf("error insert into categories")
	}

	return &product.ReturnNewCategory{
		Id:           category_id,
		CategoryName: category_name.CategoryName,
	}, nil
}

func (s *Server) ChangeProduct(ctx context.Context, product *product.ReturnNewProduct) (*product.ReturnNewProduct, error) {
	query := "UPDATE products SET product_name = $1, category_id = $2, price = $3, image_url = $4, availability_of_pieces = $5 WHERE id = $6"

	rows, err := s.DB.Exec(query, product.ProductName, product.CategoryId, product.Price, product.ImageUrl, product.AvailabilityOfPieces, product.Id)
	if err != nil {
		return nil, fmt.Errorf("error update product")
	}

	rowUpdated, _ := rows.RowsAffected()

	if rowUpdated != 1 {
		return nil, fmt.Errorf("product is not defined")
	}

	return product, nil
}

func (s *Server) ChangeCategory(ctx context.Context, category *product.ReturnNewCategory) (*product.ReturnNewCategory, error) {
	query := "UPDATE categories SET category_name = $1 WHERE id = $2"

	rows, err := s.DB.Exec(query, category.CategoryName, category.Id)
	if err != nil {
		return nil, fmt.Errorf("error update category")
	}

	rowUpdated, _ := rows.RowsAffected()

	if rowUpdated != 1 {
		return nil, fmt.Errorf("category is not defined")
	}

	return category, nil
}

func (s *Server) RemoveProduct(ctx context.Context, productId *product.ProductId) (*product.RemoveMessage, error) {
	query := "DELETE FROM products WHERE id = $1"

	rows, err := s.DB.Exec(query, productId.Id)
	if err != nil {
		return nil, fmt.Errorf("error remove product")
	}

	rowsRemove, _ := rows.RowsAffected()
	if rowsRemove != 1 {
		return nil, fmt.Errorf("product is not defined")
	}

	return &product.RemoveMessage{
		Id:      productId.Id,
		Status:  true,
		Message: "product remove successfully",
	}, nil
}

func (s *Server) RemoveCategory(ctx context.Context, categoryId *product.CategoryId) (*product.RemoveMessage, error) {
	query := "DELETE FROM categories WHERE id = $1"

	rows, err := s.DB.Exec(query, categoryId.Id)
	if err != nil {
		return nil, fmt.Errorf("error remove category")
	}

	rowsRemove, _ := rows.RowsAffected()
	if rowsRemove != 1 {
		return nil, fmt.Errorf("category is not defined")
	}

	return &product.RemoveMessage{
		Id:      categoryId.Id,
		Status:  true,
		Message: "category remove successfully",
	}, nil
}

func (s *Server) GetProductName(ctx context.Context, product_name *product.ProductName) (*product.ReturnResult, error) {
	query := "SELECT product_name FROM products WHERE product_name = $1"

	var productName string

	err := s.DB.QueryRowContext(ctx, query, product_name.ProductName).Scan(&productName)

	if err == sql.ErrNoRows {
		return &product.ReturnResult{Result: false}, nil
	}

	if err != nil {
		return nil, fmt.Errorf("%v", err.Error())
	}

	return &product.ReturnResult{Result: true}, nil
}

func (s *Server) CreateSubcategory(ctx context.Context, newSubcategory *product.NewSubcategory) (*product.ReturnNewSubcategory, error) {
	query := "INSERT INTO subcategories (category_name, category_id) VALUES ($1, $2) RETURNING id"

	var category_id int64

	err := s.DB.QueryRow(query, newSubcategory.CategoryName, newSubcategory.CategoryId).Scan(&category_id)
	if err != nil {
		return nil, fmt.Errorf("error insert into subcategories")
	}

	return &product.ReturnNewSubcategory{
		Id:           category_id,
		CategoryName: newSubcategory.CategoryName,
	}, nil
}

func (s *Server) ChangeSubcategory(ctx context.Context, subcategory *product.ReturnNewSubcategory) (*product.ReturnNewSubcategory, error) {
	query := "UPDATE subcategories SET category_name = $1, category_id = $2 WHERE id = $3"

	rows, err := s.DB.Exec(query, subcategory.CategoryName, subcategory.CategoryId, subcategory.Id)
	if err != nil {
		return nil, fmt.Errorf("error update subcategories")
	}

	rowUpdated, _ := rows.RowsAffected()

	if rowUpdated != 1 {
		return nil, fmt.Errorf("subcategory is not defined")
	}

	return subcategory, nil
}

func (s *Server) RemoveSubcategory(ctx context.Context, categoryId *product.CategoryId) (*product.RemoveMessage, error) {
	query := "DELETE FROM subcategories WHERE id = $1"

	rows, err := s.DB.Exec(query, categoryId.Id)
	if err != nil {
		return nil, fmt.Errorf("error remove category")
	}

	rowsRemove, _ := rows.RowsAffected()
	if rowsRemove != 1 {
		return nil, fmt.Errorf("category is not defined")
	}

	return &product.RemoveMessage{
		Id:      categoryId.Id,
		Status:  true,
		Message: "category remove successfully",
	}, nil
}
