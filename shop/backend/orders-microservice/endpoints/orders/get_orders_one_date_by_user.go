package orders

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"orders-microservice/helpers"
	"orders-microservice/models"
	"time"
)

func GetOrdersOneDateByUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var emailParam = r.URL.Query().Get("email")
	var phoneParam = r.URL.Query().Get("phone")
	var userIdParam = r.URL.Query().Get("user_id")
	var dateParam = r.URL.Query().Get("date")
	var orders []models.FullOrder
	var ordersMap = make(map[int]*models.FullOrder)
	var ctx = r.Context()
	if dateParam == "" {
		http.Error(w, "date can't be empty", http.StatusBadRequest)
		return
	}
	fmt.Println(dateParam)
	_, err := time.Parse("2006-01-02", dateParam)
	if err != nil {
		http.Error(w, "invalid date format, use YYYY-MM-DD", http.StatusBadRequest)
		return
	}
	var query string = `SELECT orders.id as order_id, user_id, email, phone, status, total_price, created_at, order_items.id AS order_item_id, 
						product_id, quantity, order_items.price, product_name, category_id, category_name, image_url
						FROM orders 
						JOIN order_items ON orders.id = order_items.order_id
						JOIN products ON order_items.product_id = products.id
						JOIN categories ON products.category_id = categories.id`
	if emailParam != "" {

		rows, err := db.QueryContext(ctx, query+" "+helpers.SqlQueryWithParamAndDate("email"), emailParam, dateParam)

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
			orders = append(orders, *order)
		}
	}

	if phoneParam != "" {

		rows, err := db.QueryContext(ctx, query+" "+helpers.SqlQueryWithParamAndDate("phone"), phoneParam, dateParam)

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
			orders = append(orders, *order)
		}

	}

	if userIdParam != "" {

		rows, err := db.QueryContext(ctx, query+" "+helpers.SqlQueryWithParamAndDate("user_id"), userIdParam, dateParam)

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
			orders = append(orders, *order)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}
