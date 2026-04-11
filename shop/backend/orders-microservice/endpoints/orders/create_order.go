package orders

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"orders-microservice/models"
)

func CreateOrder(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var order models.NewOrder

	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if (order.Email == nil || *order.Email == "") && (order.Phone == nil || *order.Phone == "") && (order.UserId == nil || *order.UserId <= 0) {
		http.Error(w, "Contact info dont't empty", http.StatusBadRequest)
		return
	}

	if len(order.OrderItems) == 0 {
		http.Error(w, "Order must contain items", http.StatusBadRequest)
		return
	}

	if order.Status == nil || *order.Status != "pending" {
		val := "pending"
		order.Status = &val
	}

	ctx := r.Context()

	queryOrder := `
		INSERT INTO orders (user_id, email, phone, status, total_price)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	queryOrderItems := `
		INSERT INTO order_items (order_id, product_id, quantity, price)
		VALUES ($1, $2, $3, $4)
	`

	var productsIds []int
	var productsCheck []models.ProductsCheck
	var productsFromOrder []models.OrderItem
	var problemProducts []models.ProblemProducts
	for _, product := range order.OrderItems {
		productsIds = append(productsIds, product.ProductId)
		productsFromOrder = append(productsFromOrder, product)
	}

	var queryCheckQuantityProduct = `SELECT id, availability_of_pieces, product_name FROM products WHERE id = ANY($1)`

	rows, err := db.QueryContext(ctx, queryCheckQuantityProduct, productsIds)

	if err != nil {
		http.Error(w, "error check products for stock", http.StatusInternalServerError)
		return
	}

	for rows.Next() {
		var productCheck models.ProductsCheck
		err := rows.Scan(&productCheck.ProductId, &productCheck.AvailabilityOfPieces, &productCheck.ProductName)
		if err != nil {
			http.Error(w, "Error reading row", http.StatusInternalServerError)
			return
		}

		productsCheck = append(productsCheck, productCheck)
	}

	for i, product := range productsCheck {
		if product.AvailabilityOfPieces < productsFromOrder[i].Quantity {
			problemProducts = append(problemProducts, models.ProblemProducts(product))
		}
	}

	if len(problemProducts) > 0 {
		problemProductsMsg := map[string]interface{}{
			"message":          "Sorry, there were no products in stock",
			"problem_products": problemProducts,
		}
		json.NewEncoder(w).Encode(problemProductsMsg)
		return
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer func() {
		_ = tx.Rollback()
	}()

	var newOrderId models.NewOrderId
	totalPrice := 0.0

	err = tx.QueryRowContext(ctx, queryOrder,
		order.UserId,
		order.Email,
		order.Phone,
		order.Status,
		0,
	).Scan(&newOrderId.Id)

	if err != nil {
		http.Error(w, "Error creating order", http.StatusInternalServerError)
		return
	}

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

		rows, _ := res.RowsAffected()
		if rows == 0 {
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

		_, err = tx.ExecContext(ctx, queryOrderItems,
			newOrderId.Id,
			item.ProductId,
			item.Quantity,
			price,
		)

		if err != nil {
			http.Error(w, "Error creating order_items", http.StatusInternalServerError)
			return
		}

		totalPrice += price * float64(item.Quantity)
	}

	_, err = tx.ExecContext(ctx, `
		UPDATE orders SET total_price = $1 WHERE id = $2
	`, totalPrice, newOrderId.Id)

	if err != nil {
		http.Error(w, "Failed to update total price", http.StatusInternalServerError)
		return
	}

	var query string = `SELECT orders.id as order_id, user_id, email, phone, status, total_price, created_at, order_items.id AS order_item_id, 
							product_id, quantity, order_items.price, product_name, category_id, category_name, image_url
							FROM orders 
							JOIN order_items ON orders.id = order_items.order_id
							JOIN products ON order_items.product_id = products.id
							JOIN categories ON products.category_id = categories.id`

	var emailInOrder = order.Email
	var phoneInOrder = order.Phone
	var products []models.Products
	var fullOrder models.FullOrder
	if emailInOrder != nil && *emailInOrder != "" {
		query += " WHERE email = $1 AND order_id = $2"
		rows, err := tx.QueryContext(ctx, query, emailInOrder, newOrderId.Id)
		if err != nil {
			http.Error(w, "order receiving error", http.StatusInternalServerError)
			return
		}

		for rows.Next() {
			var item models.Products
			err := rows.Scan(&fullOrder.OrderId, &fullOrder.UserId, &fullOrder.Email, &fullOrder.Phone, &fullOrder.Status,
				&fullOrder.TotalPrice, &item.OrderItemId, &item.ProductId, &item.Quantity, &item.Price, &item.ProductName, &item.CategoryId, &item.CategoryName, &item.ImageUrl)

			if err != nil {
				http.Error(w, "Error reading row", http.StatusInternalServerError)
				return
			}
			products = append(products, item)
		}
	}

	if phoneInOrder != nil && *phoneInOrder != "" {
		query += " WHERE phone = $1 AND order_id = $2"
		rows, err := tx.QueryContext(ctx, query, phoneInOrder, newOrderId.Id)
		if err != nil {
			http.Error(w, "order receiving error", http.StatusInternalServerError)
			return
		}

		for rows.Next() {
			var item models.Products
			err := rows.Scan(&fullOrder.OrderId, &fullOrder.UserId, &fullOrder.Email, &fullOrder.Phone, &fullOrder.Status,
				&fullOrder.TotalPrice, &fullOrder.CreatedAt, &item.OrderItemId, &item.ProductId, &item.Quantity, &item.Price, &item.ProductName, &item.CategoryId, &item.CategoryName, &item.ImageUrl)

			if err != nil {
				http.Error(w, "Error reading row", http.StatusInternalServerError)
				return
			}
			products = append(products, item)
		}
	}

	fullOrder.Products = products

	if err := tx.Commit(); err != nil {
		http.Error(w, "Commit failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fullOrder)
}
