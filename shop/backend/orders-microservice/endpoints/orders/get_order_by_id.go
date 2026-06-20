package orders

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"orders-microservice/models"
	"strconv"

	product "github.com/AndreyLebedev1998/shop-gRPC-product"
)

func GetOrderById(w http.ResponseWriter, r *http.Request, db *sql.DB, client product.ProductsServiceClient) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var orderID = r.URL.Query().Get("order_id")
	if orderID == "" {
		http.Error(w, "order_id is not defined", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(orderID)
	if err != nil || id < 0 {
		http.Error(w, "Invalid order_id", http.StatusBadRequest)
		return
	}
	var fullOrder models.FullOrder
	var products []models.Products
	ctx := r.Context()

	var query string = `SELECT orders.id as order_id, user_id, email, phone, status, total_price, created_at, 
						product_id, quantity, order_items.price
						FROM orders 
						JOIN order_items ON orders.id = order_items.order_id
						WHERE orders.id = $1`

	rows, err := db.QueryContext(ctx, query, id)

	if err == sql.ErrNoRows {
		http.Error(w, "order is not defined", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	for rows.Next() {
		var product models.Products

		if err := rows.Scan(&fullOrder.OrderId, &fullOrder.UserId, &fullOrder.Email, &fullOrder.Phone,
			&fullOrder.Status, &fullOrder.TotalPrice, &fullOrder.CreatedAt,
			&product.ProductId, &product.Quantity, &product.Price); err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		fullOrder.Products = append(fullOrder.Products, product)
	}

	fmt.Println(fullOrder)

	var productsIds []int64

	for _, p := range fullOrder.Products {
		productsIds = append(productsIds, int64(p.ProductId))
	}

	fmt.Println(productsIds)

	resp, err := client.GetProductsByIds(ctx, &product.GetProductsRequest{
		ProductIds: productsIds,
	})

	if err != nil {
		fmt.Println(err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	if resp != nil {
		for _, p := range resp.Products {
			products = append(products, models.Products{
				ProductId:    int(p.Id),
				ProductName:  p.ProductName,
				CategoryId:   int(p.CategoryId),
				CategoryName: p.CategoryName,
				ImageUrl:     &p.ImageUrl,
			})
		}
	}

	for i, product := range fullOrder.Products {
		for _, p := range products {
			if product.ProductId == p.ProductId {
				fullOrder.Products[i].ProductName = p.ProductName
				fullOrder.Products[i].CategoryName = p.CategoryName
				fullOrder.Products[i].CategoryId = p.CategoryId
				fullOrder.Products[i].ImageUrl = p.ImageUrl
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fullOrder)
}
