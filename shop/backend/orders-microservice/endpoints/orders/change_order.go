package orders

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"orders-microservice/helpers"
	"orders-microservice/models"

	product "github.com/AndreyLebedev1998/shop-gRPC-product"
)

func ChangeOrder(w http.ResponseWriter, r *http.Request, db *sql.DB, client product.ProductsServiceClient) {
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

	var oldItems []models.OrderItem

	rows, err := tx.QueryContext(ctx, oldItemsQuery, order.Id)
	if err != nil {
		http.Error(w, "Error fetching old items", http.StatusInternalServerError)
		return
	}

	for rows.Next() {
		var item models.OrderItem
		if err := rows.Scan(&item.ProductId, &item.Quantity); err != nil {
			rows.Close()
			http.Error(w, "Scan error", http.StatusInternalServerError)
			return
		}
		oldItems = append(oldItems, item)
	}
	defer rows.Close()

	var deleteQuery string = `DELETE FROM order_items WHERE order_id = $1`

	resDel, err := tx.ExecContext(ctx, deleteQuery, order.Id)

	if err != nil {
		http.Error(w, "Error delete order_items", http.StatusInternalServerError)
		return
	}

	rowsDeleted, err := resDel.RowsAffected()

	if rowsDeleted != int64(len(order.OrderItems)) {
		http.Error(w, "order_id is not defined", http.StatusInternalServerError)
		return
	}

	var insertQuery string = `INSERT INTO order_items (order_id, product_id, quantity)
									VALUES ($1, $2, $3)`

	for _, item := range order.OrderItems {
		resIns, err := tx.ExecContext(ctx, insertQuery, order.Id, item.ProductId, item.Quantity)

		if err != nil {
			http.Error(w, "error INSERT INTO order_items 1", http.StatusInternalServerError)
			return
		}

		rowsUpdated, _ := resIns.RowsAffected()

		if rowsUpdated != int64(1) {
			http.Error(w, "error INSERT INTO order_items 2", http.StatusInternalServerError)
			return
		}
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

	var problemProducts []models.ProblemProducts

	for i, product := range resp.Products {
		if product.AvailabilityOfPieces < int64(order.OrderItems[i].Quantity) {
			var problemProduct models.ProblemProducts
			problemProduct.AvailabilityOfPieces = int(product.AvailabilityOfPieces)
			problemProduct.ImageUrl = &product.ImageUrl
			problemProduct.ProductId = int(product.Id)
			problemProduct.ProductName = product.ProductName
			problemProducts = append(problemProducts, models.ProblemProducts(problemProduct))
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

	var fullOrder models.FullOrder
	var productsInOrder []models.Products

	for _, product := range resp.Products {
		for _, p := range order.OrderItems {
			if product.Id == int64(p.ProductId) {
				var oneProduct models.Products
				oneProduct.ProductId = int(product.Id)
				oneProduct.Quantity = p.Quantity
				oneProduct.Price = product.Price
				oneProduct.ProductName = product.ProductName
				oneProduct.CategoryId = int(product.CategoryId)
				oneProduct.CategoryName = product.CategoryName
				oneProduct.ImageUrl = &product.ImageUrl

				productsInOrder = append(productsInOrder, oneProduct)
			}
		}
	}

	fullOrder.OrderId = order.Id
	fullOrder.UserId = order.UserId
	fullOrder.Email = order.Email
	fullOrder.Phone = order.Phone
	fullOrder.Status = *order.Status
	fullOrder.TotalPrice = order.TotalPrice
	fullOrder.CreatedAt = order.CreatedAt
	fullOrder.UpdatedAt = order.UpdatedAt
	fullOrder.Products = productsInOrder

	grpcNewItems := make([]*product.UpdateProductQuantity, 0, len(oldItems))
	grpcOldItems := make([]*product.UpdateProductQuantity, 0, len(oldItems))

	for _, p := range order.OrderItems {
		grpcNewItems = append(grpcNewItems, &product.UpdateProductQuantity{
			ProductId: int64(p.ProductId),
			Quantity:  int64(p.Quantity),
		})
	}

	for _, p := range oldItems {
		grpcOldItems = append(grpcOldItems, &product.UpdateProductQuantity{
			ProductId: int64(p.ProductId),
			Quantity:  int64(p.Quantity),
		})
	}

	response, err := client.UpdateProductQuantityByIds(ctx, &product.UpdateProductQuantityRequest{
		NewItems: grpcNewItems,
		OldItems: grpcOldItems,
		IsCreate: false,
	})

	fmt.Println(response)

	if err != nil {
		http.Error(w, "Server error 1", http.StatusInternalServerError)
		return
	}

	if !response.Success {
		http.Error(w, "Server error 2", http.StatusInternalServerError)
		return
	}

	totalPrice := 0.0
	for _, p := range productsInOrder {
		for _, product := range order.OrderItems {
			if int(p.ProductId) == int(product.ProductId) {
				totalPrice += (p.Price * float64(product.Quantity))
				continue
			}
		}
	}

	fullOrder.TotalPrice = totalPrice

	if err := tx.Commit(); err != nil {
		http.Error(w, "Commit failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fullOrder)
}
