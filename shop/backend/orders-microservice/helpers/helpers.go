package helpers

import (
	"database/sql"
	"fmt"
	"orders-microservice/models"
)

func SqlQueryWithParam(columnName string) string {
	return fmt.Sprintf("WHERE %s = $1", columnName)
}

func SqlQueryWithParamAndDate(columnName string) string {
	return fmt.Sprintf("WHERE %s = $1 AND DATE(created_at) = $2", columnName)
}

func ForRowsAfterQuery(rows *sql.Rows, ordersMap map[int]*models.FullOrder) error {
	for rows.Next() {
		var fullOrder models.FullOrder
		var product models.Products
		if err := rows.Scan(&fullOrder.OrderId, &fullOrder.UserId, &fullOrder.Email, &fullOrder.Phone,
			&fullOrder.Status, &fullOrder.TotalPrice, &fullOrder.CreatedAt,
			&product.OrderItemId, &product.ProductId, &product.Quantity, &product.Price,
			&product.ProductName, &product.CategoryId, &product.CategoryName, &product.ImageUrl); err != nil {
			return fmt.Errorf("error scanning row: %w", err)
		}

		if o, ok := ordersMap[fullOrder.OrderId]; !ok {
			fullOrder.Products = []models.Products{product}
			ordersMap[fullOrder.OrderId] = &fullOrder
		} else {
			o.Products = append(o.Products, product)
		}
	}
	return rows.Err()
}
