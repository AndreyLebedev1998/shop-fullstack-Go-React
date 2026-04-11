package orders

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"orders-microservice/models"
)

func ChangeOrder(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var order models.Order

	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if len(order.OrderItems) == 0 {
		http.Error(w, "Order must contain items", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer func() {
		_ = tx.Rollback()
	}()

	oldItemsQuery := `
		SELECT product_id, quantity
		FROM order_items
		WHERE order_id = $1
	`

	rows, err := tx.QueryContext(ctx, oldItemsQuery, order.Id)
	if err != nil {
		http.Error(w, "Error fetching old items", http.StatusInternalServerError)
		return
	}

	var oldItems []models.OrderItem

	for rows.Next() {
		var item models.OrderItem
		if err := rows.Scan(&item.ProductId, &item.Quantity); err != nil {
			rows.Close()
			http.Error(w, "Scan error", http.StatusInternalServerError)
			return
		}
		oldItems = append(oldItems, item)
	}
	rows.Close()

	for _, item := range oldItems {
		_, err := tx.ExecContext(ctx, `
			UPDATE products
			SET availability_of_pieces = availability_of_pieces + $1
			WHERE id = $2
		`, item.Quantity, item.ProductId)

		if err != nil {
			http.Error(w, "Failed to restore stock", http.StatusInternalServerError)
			return
		}
	}

	_, err = tx.ExecContext(ctx, `
		DELETE FROM order_items WHERE order_id = $1
	`, order.Id)

	if err != nil {
		http.Error(w, "Delete error", http.StatusInternalServerError)
		return
	}

	var totalPrice float64

	for _, item := range order.OrderItems {

		res, err := tx.ExecContext(ctx, `
			UPDATE products
			SET availability_of_pieces = availability_of_pieces - $1
			WHERE id = $2 AND availability_of_pieces >= $1
		`, item.Quantity, item.ProductId)

		if err != nil {
			http.Error(w, "Stock update failed", http.StatusInternalServerError)
			return
		}

		rowsAffected, _ := res.RowsAffected()
		if rowsAffected == 0 {
			http.Error(w, "Not enough stock", http.StatusBadRequest)
			return
		}

		var price float64
		err = tx.QueryRowContext(ctx, `
			SELECT price FROM products WHERE id = $1
		`, item.ProductId).Scan(&price)

		if err != nil {
			http.Error(w, "Product not found", http.StatusBadRequest)
			return
		}

		_, err = tx.ExecContext(ctx, `
			INSERT INTO order_items (order_id, product_id, quantity, price)
			VALUES ($1, $2, $3, $4)
		`, order.Id, item.ProductId, item.Quantity, price)

		if err != nil {
			http.Error(w, "Insert item failed", http.StatusInternalServerError)
			return
		}

		totalPrice += price * float64(item.Quantity)
	}

	_, err = tx.ExecContext(ctx, `
		UPDATE orders
		SET total_price = $1
		WHERE id = $2
	`, totalPrice, order.Id)

	if err != nil {
		http.Error(w, "Order update failed", http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, "Commit failed", http.StatusInternalServerError)
		return
	}

	var query = `
	SELECT 
		orders.id as order_id,
		user_id,
		email,
		phone,
		status,
		total_price,
		created_at,
		order_items.id AS order_item_id,
		product_id,
		quantity,
		order_items.price,
		product_name,
		category_id,
		category_name,
		image_url
	FROM orders 
	JOIN order_items ON orders.id = order_items.order_id
	JOIN products ON order_items.product_id = products.id
	JOIN categories ON products.category_id = categories.id
	WHERE orders.id = $1
	`

	rows, err = db.QueryContext(ctx, query, order.Id)
	if err != nil {
		http.Error(w, "order receiving error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var fullOrder models.FullOrder
	var products []models.Products

	for rows.Next() {
		var item models.Products

		err := rows.Scan(
			&fullOrder.OrderId,
			&fullOrder.UserId,
			&fullOrder.Email,
			&fullOrder.Phone,
			&fullOrder.Status,
			&fullOrder.TotalPrice,
			&fullOrder.CreatedAt,
			&item.OrderItemId,
			&item.ProductId,
			&item.Quantity,
			&item.Price,
			&item.ProductName,
			&item.CategoryId,
			&item.CategoryName,
			&item.ImageUrl,
		)

		if err != nil {
			http.Error(w, "Error reading row", http.StatusInternalServerError)
			return
		}

		products = append(products, item)
	}

	fullOrder.Products = products

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fullOrder)
}
