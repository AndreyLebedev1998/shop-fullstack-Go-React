package orders

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"orders-microservice/helpers"
	"orders-microservice/models"

	product "github.com/AndreyLebedev1998/shop-gRPC-product"
)

func markOrderFailed(ctx context.Context, db *sql.DB, orderId int) {
	_, _ = db.ExecContext(ctx, "UPDATE orders SET status = $1 WHERE id = $2", "failed", orderId)
}

func CreateOrder(w http.ResponseWriter, r *http.Request, db *sql.DB, client product.ProductsServiceClient) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var order models.NewOrder
	ctx := r.Context()

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

	var product_ids []int

	for _, p := range order.OrderItems {
		product_ids = append(product_ids, p.ProductId)
	}

	resp, err := client.GetProductsByIds(ctx, &product.GetProductsRequest{
		ProductIds: helpers.ConvertIntToInt64(product_ids),
	})

	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	if order.Status == nil || *order.Status != "pending" {
		val := "pending"
		order.Status = &val
	}

	queryOrder := `
		INSERT INTO orders (user_id, email, phone, status, total_price)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	queryOrderItems := `
		INSERT INTO order_items (order_id, product_id, quantity, price)
		VALUES ($1, $2, $3, $4)
	`

	var productsFromOrder []models.OrderItem
	var problemProducts []models.ProblemProducts
	var productsInOrder []models.ProductsForOrder
	for _, product := range order.OrderItems {
		productsFromOrder = append(productsFromOrder, product)
	}

	for i, product := range resp.Products {
		if product.AvailabilityOfPieces < int64(productsFromOrder[i].Quantity) {
			var problemProduct models.ProblemProducts
			problemProduct.AvailabilityOfPieces = int(product.AvailabilityOfPieces)
			problemProduct.ImageUrl = &product.ImageUrl
			problemProduct.ProductId = int(product.Id)
			problemProduct.ProductName = product.ProductName
			problemProducts = append(problemProducts, models.ProblemProducts(problemProduct))
		}

		var productInOrder models.ProductsForOrder

		productInOrder.Id = int(product.Id)
		productInOrder.ProductName = product.ProductName
		productInOrder.CategoryId = int(product.CategoryId)
		productInOrder.CategoryName = product.CategoryName
		productInOrder.ImageUrl = &product.ImageUrl
		productInOrder.Quantity = int(product.AvailabilityOfPieces)
		productInOrder.Price = product.Price

		productsInOrder = append(productsInOrder, productInOrder)
	}

	if len(problemProducts) > 0 {
		problemProductsMsg := map[string]interface{}{
			"message":          "Sorry, there were no products in stock",
			"problem_products": problemProducts,
		}
		json.NewEncoder(w).Encode(problemProductsMsg)
		return
	}

	totalPrice := 0.0
	for _, p := range productsInOrder {
		for _, product := range order.OrderItems {
			if int(p.Id) == int(product.ProductId) {
				totalPrice += (p.Price * float64(product.Quantity))
				continue
			}
		}
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

	_, err = tx.ExecContext(ctx, `
		UPDATE orders SET total_price = $1 WHERE id = $2
	`, totalPrice, newOrderId.Id)

	if err != nil {
		http.Error(w, "Failed to update total price", http.StatusInternalServerError)
		return
	}

	for _, item := range order.OrderItems {
		for _, p := range productsInOrder {
			if item.ProductId == p.Id {
				price := p.Price * float64(item.Quantity)
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
			}
		}
	}

	var query string = `SELECT orders.id as order_id, user_id, email, phone, status, total_price, created_at,
							product_id, quantity, price
							FROM orders
							JOIN order_items ON orders.id = order_items.order_id`

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
				&fullOrder.TotalPrice, &fullOrder.CreatedAt, &item.ProductId, &item.Quantity, &item.Price)

			if err != nil {
				http.Error(w, "Error reading row", http.StatusInternalServerError)
				return
			}
			products = append(products, item)
		}

		for i := range products {
			for _, p := range productsInOrder {
				if products[i].ProductId == p.Id {
					products[i].CategoryId = p.CategoryId
					products[i].CategoryName = p.CategoryName
					products[i].ImageUrl = p.ImageUrl
					products[i].ProductName = p.ProductName
				}
			}
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
				&fullOrder.TotalPrice, &fullOrder.CreatedAt, &item.ProductId, &item.Quantity, &item.Price)

			if err != nil {
				http.Error(w, "Error reading row", http.StatusInternalServerError)
				return
			}
			products = append(products, item)
		}

		for i := range products {
			for _, p := range productsInOrder {
				if products[i].ProductId == p.Id {
					products[i].CategoryId = p.CategoryId
					products[i].CategoryName = p.CategoryName
					products[i].ImageUrl = p.ImageUrl
					products[i].ProductName = p.ProductName
				}
			}
		}
	}

	fullOrder.Products = products

	if err := tx.Commit(); err != nil {
		http.Error(w, "Commit failed", http.StatusInternalServerError)
		return
	}

	grpcItems := make([]*product.UpdateProductQuantity, 0, len(productsFromOrder))

	for _, p := range productsFromOrder {
		grpcItems = append(grpcItems, &product.UpdateProductQuantity{
			ProductId: int64(p.ProductId),
			Quantity:  int64(p.Quantity),
		})
	}
	response, err := client.UpdateProductQuantityByIds(ctx, &product.UpdateProductQuantityRequest{
		NewItems: grpcItems,
		OldItems: grpcItems,
		IsCreate: true,
	})

	if err != nil {
		markOrderFailed(ctx, db, newOrderId.Id)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !response.Success {
		markOrderFailed(ctx, db, newOrderId.Id)
		http.Error(w, response.Message, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fullOrder)
}
