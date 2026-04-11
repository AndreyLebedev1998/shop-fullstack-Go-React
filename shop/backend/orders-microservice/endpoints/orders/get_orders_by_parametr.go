package orders

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"orders-microservice/helpers"
	"orders-microservice/models"
)

func GetOrdersByParametr(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var emailParametr string = r.URL.Query().Get("email")
	var phoneParamert string = r.URL.Query().Get("phone")
	var userIdParametr string = r.URL.Query().Get("user_id")
	var ctx = r.Context()
	var fullOrders []models.FullOrder
	ordersMap := make(map[int]*models.FullOrder)
	var query string = `SELECT orders.id as order_id, user_id, email, phone, status, total_price, created_at, order_items.id AS order_item_id, 
							product_id, quantity, order_items.price, product_name, category_id, category_name, image_url
							FROM orders 
							JOIN order_items ON orders.id = order_items.order_id
							JOIN products ON order_items.product_id = products.id
							JOIN categories ON products.category_id = categories.id`

	if emailParametr != "" {

		rows, err := db.QueryContext(ctx, query+" "+helpers.SqlQueryWithParam("email"), emailParametr)

		if err != nil {
			http.Error(w, "Error while querying the database", http.StatusInternalServerError)
			return
		}

		defer rows.Close()

		if err := helpers.ForRowsAfterQuery(rows, ordersMap); err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		for _, order := range ordersMap {
			fullOrders = append(fullOrders, *order)
		}
	}

	if phoneParamert != "" {
		rows, err := db.QueryContext(ctx, query+" "+helpers.SqlQueryWithParam("email"), emailParametr)

		if err != nil {
			http.Error(w, "Error while querying the database", http.StatusInternalServerError)
			return
		}

		defer rows.Close()

		if err := helpers.ForRowsAfterQuery(rows, ordersMap); err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		for _, order := range ordersMap {
			fullOrders = append(fullOrders, *order)
		}
	}

	if userIdParametr != "" {
		rows, err := db.QueryContext(ctx, query+" "+helpers.SqlQueryWithParam("email"), emailParametr)

		if err != nil {
			http.Error(w, "Error while querying the database", http.StatusInternalServerError)
			return
		}

		defer rows.Close()

		if err := helpers.ForRowsAfterQuery(rows, ordersMap); err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		for _, order := range ordersMap {
			fullOrders = append(fullOrders, *order)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fullOrders)
}
