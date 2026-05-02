package remove

import (
	"encoding/json"
	"net/http"
	"strconv"

	product "github.com/AndreyLebedev1998/shop-gRPC-product"
)

func RemoveCategory(w http.ResponseWriter, r *http.Request, client product.ProductsServiceClient) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var ctx = r.Context()

	categoryIdStr := r.URL.Query().Get("id")
	categoryId, err := strconv.Atoi(categoryIdStr)
	if err != nil || categoryId <= 0 {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		return
	}

	resp, err := client.RemoveCategory(ctx, &product.CategoryId{
		Id: int64(categoryId),
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application-json")
	json.NewEncoder(w).Encode(resp)
}
