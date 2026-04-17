package grpcFunctionsQuery

import (
	"context"
	"database/sql"
	"fmt"
	"products-microservice/models"

	"github.com/lib/pq"
)

func GetProductByIds(db *sql.DB, products_ids []int) ([]models.ProductWithCategory, error) {
	query := `SELECT products.id, product_name, category_id, price, image_url, availability_of_pieces, category_name FROM products
			  JOIN categories ON products.category_id = categories.id
			  WHERE products.id = ANY($1)`
	var prodcuts []models.ProductWithCategory

	rows, err := db.Query(query, pq.Array(products_ids))

	if err != nil {
		return nil, fmt.Errorf("error database")
	}

	for rows.Next() {
		var product models.ProductWithCategory
		err := rows.Scan(&product.Id, &product.ProductName, &product.CategoryId, &product.Price, &product.ImageUrl, &product.AvailabilityOfPieces, &product.CategoryName)
		if err != nil {
			return nil, fmt.Errorf("error reading row")
		}
		prodcuts = append(prodcuts, product)
	}

	return prodcuts, nil
}

func UpdateProductsQuantityByIds(db *sql.DB, products []models.UpdateProductForGRPC) (models.MessageUpdatedQuantityProducts, error) {
	query := "UPDATE products SET availability_of_pieces = availability_of_pieces - $1 WHERE id = $2 AND availability_of_pieces >= $1"

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		response := models.MessageUpdatedQuantityProducts{
			Success: false,
			Message: "transaction error",
		}
		return response, fmt.Errorf("Server error")
	}

	defer func() {
		_ = tx.Rollback()
	}()

	for _, p := range products {
		res, err := tx.Exec(query, p.Quantity, p.ProductId)
		if err != nil {
			response := models.MessageUpdatedQuantityProducts{
				Success: false,
				Message: "Server error",
			}
			return response, fmt.Errorf("Server error")
		}

		rowsUpdated, _ := res.RowsAffected()
		if rowsUpdated != 1 {
			response := models.MessageUpdatedQuantityProducts{
				Success: false,
				Message: "Updated products quantity error",
			}
			return response, fmt.Errorf("Server error")
		}
	}

	response := models.MessageUpdatedQuantityProducts{
		Success: true,
		Message: "Products in the warehouse have been successfully updated.",
	}

	if err := tx.Commit(); err != nil {
		response := models.MessageUpdatedQuantityProducts{
			Success: false,
			Message: "Commit failed",
		}
		return response, fmt.Errorf("Server error")
	}

	return response, nil
}
