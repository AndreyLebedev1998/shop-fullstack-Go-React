package change

import (
	"admin-microservice/models"
	"encoding/json"
	"net/http"
	"strconv"

	product "github.com/AndreyLebedev1998/shop-gRPC-product"
)

func ChangeCategory(w http.ResponseWriter, r *http.Request, client product.ProductsServiceClient) {
	if r.Method != http.MethodPatch {
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

	categoryIdStr := r.FormValue("id")
	categoryId, err := strconv.Atoi(categoryIdStr)
	if err != nil || categoryId <= 0 {
		http.Error(w, "Invalid category id", http.StatusBadRequest)
		return
	}

	var grpcNewCategory = &product.ReturnNewCategory{
		Id:           int64(categoryId),
		CategoryName: categoryName.CategoryName,
	}

	resp, err := client.ChangeCategory(ctx, grpcNewCategory)

	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
