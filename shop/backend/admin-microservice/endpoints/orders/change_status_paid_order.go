package orders

import (
	"admin-microservice/constants"
	"admin-microservice/helpers"
	"admin-microservice/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func ChangeStatusPaidOrder(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var orderId = r.URL.Query().Get("order_id")
	var orderStatusPaid models.OrderStatusPaid
	var ctx = r.Context()

	if orderId == "" {
		http.Error(w, "order_id is not defined", http.StatusBadRequest)
		return
	}

	orderIdNum, err := strconv.Atoi(orderId)
	if err != nil || orderIdNum < 0 {
		http.Error(w, "order_id is not valid", http.StatusBadGateway)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&orderStatusPaid); err != nil {
		http.Error(w, "status_paid is not defined", http.StatusBadRequest)
		return
	}

	var statusStr constants.OrderStatusPaid = constants.OrderStatusPaid(orderStatusPaid.StatusPaid)

	isValid := helpers.IsValidStatusPaid(statusStr)

	fmt.Println(isValid)

	if isValid {
		query := `UPDATE orders SET status_paid = $1 WHERE id = $2`

		res, err := db.ExecContext(ctx, query, orderStatusPaid.StatusPaid, orderId)

		if err != nil {
			http.Error(w, "Error update order", http.StatusInternalServerError)
			return
		}

		rowsUpdated, err := res.RowsAffected()

		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		if rowsUpdated == 0 {
			http.Error(w, "order is not defined", http.StatusBadRequest)
			return
		}

		var success = map[string]string{
			"response": "Status changed successfully",
			"status":   string(statusStr),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(success)
	} else {
		http.Error(w, "status_paid is not valid", http.StatusBadRequest)
		return
	}
}
