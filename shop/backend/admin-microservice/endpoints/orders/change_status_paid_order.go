package orders

import (
	"admin-microservice/constants"
	"admin-microservice/helpers"
	"admin-microservice/models"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	order "github.com/AndreyLebedev1998/shop-gRPC-orders"
)

func ChangeStatusPaidOrder(w http.ResponseWriter, r *http.Request, clientOrder order.OrderServiceClient) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var status models.OrderStatusPaid
	var ctx = r.Context()
	var orderId = r.URL.Query().Get("order_id")
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

	var statusStr constants.OrderStatusPaid = constants.OrderStatusPaid(status.StatusPaid)

	isValid := helpers.IsValidStatusPaid(statusStr)

	if isValid {
		resp, err := clientOrder.ChangeOrderStatusPaid(ctx, &order.OrderStatusPaid{
			Id:         int64(id),
			StatusPaid: status.StatusPaid,
		})

		fmt.Println(err)

		if err != nil {
			http.Error(w, "Server error", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	} else {
		http.Error(w, "status is not valid", http.StatusBadRequest)
		return
	}
}
