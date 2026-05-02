package create

import (
	"admin-microservice/models"
	"encoding/json"
	"net/http"

	product "github.com/AndreyLebedev1998/shop-gRPC-product"
)

func CreateCategory(w http.ResponseWriter, r *http.Request, client product.ProductsServiceClient) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var categoryName models.NewCategory
	var ctx = r.Context()

	if err := json.NewDecoder(r.Body).Decode(&categoryName); err != nil {
		http.Error(w, "category_name is not defined", http.StatusBadRequest)
		return
	}

	if categoryName.CategoryName == "" {
		http.Error(w, "category_name can't be empty", http.StatusBadRequest)
		return
	}

	var grpcNewCategory = &product.NewCategory{
		CategoryName: categoryName.CategoryName,
	}

	resp, err := client.CreateCategory(ctx, grpcNewCategory)

	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
