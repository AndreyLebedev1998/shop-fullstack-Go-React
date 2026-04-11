package orders

import (
	"admin-microservice/constants"
	"admin-microservice/helpers"
	"admin-microservice/models"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
)

func ChangeStatusOrder(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var status models.OrderStatus
	var orderId = r.URL.Query().Get("order_id")
	var ctx = r.Context()
	if orderId == "" {
		http.Error(w, "order_id is not defined", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(orderId)
	if err != nil || id < 0 {
		http.Error(w, "Invalid order_id", http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&status); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var statusStr constants.OrderStatus = constants.OrderStatus(status.Status)

	isValid := helpers.IsValidStatus(statusStr)

	if isValid {
		var query = `UPDATE orders SET status = $1 WHERE id = $2`

		res, err := db.ExecContext(ctx, query, statusStr, orderId)
		if err != nil {
			http.Error(w, "Server Error", http.StatusInternalServerError)
			return
		}

		updatedRow, err := res.RowsAffected()

		if err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}

		if updatedRow == 0 {
			http.Error(w, "Order is not defined", http.StatusBadRequest)
			return
		}

		var success = map[string]string{
			"response": "Status changed successfully",
			"status":   string(statusStr),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(success)
	} else {
		http.Error(w, "status is not valid", http.StatusBadRequest)
		return
	}
}
