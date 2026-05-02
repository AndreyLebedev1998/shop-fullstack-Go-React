package remove

import (
	"encoding/json"
	"net/http"
	"strconv"

	product "github.com/AndreyLebedev1998/shop-gRPC-product"
)

func RemoveProduct(w http.ResponseWriter, r *http.Request, client product.ProductsServiceClient) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var ctx = r.Context()

	productIdStr := r.URL.Query().Get("id")
	productId, err := strconv.Atoi(productIdStr)
	if err != nil || productId <= 0 {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		return
	}

	resp, err := client.RemoveProduct(ctx, &product.ProductId{
		Id: int64(productId),
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application-json")
	json.NewEncoder(w).Encode(resp)
}
