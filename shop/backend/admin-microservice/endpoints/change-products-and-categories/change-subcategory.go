package change

import (
	"admin-microservice/models"
	"encoding/json"
	"net/http"
	"strconv"

	product "github.com/AndreyLebedev1998/shop-gRPC-product"
)

func ChangeSubcategory(w http.ResponseWriter, r *http.Request, client product.ProductsServiceClient) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var newSubcategory models.NewSubcategory
	var ctx = r.Context()

	if err := json.NewDecoder(r.Body).Decode(&newSubcategory); err != nil {
		http.Error(w, "category_name is not defined", http.StatusBadRequest)
		return
	}

	if newSubcategory.CategoryName == "" {
		http.Error(w, "category_name can't be empty", http.StatusBadRequest)
		return
	}

	categoryIdStr := r.FormValue("id")
	categoryId, err := strconv.Atoi(categoryIdStr)
	if err != nil || categoryId <= 0 {
		http.Error(w, "Invalid category id", http.StatusBadRequest)
		return
	}

	var grpcNewSubcategory = &product.ReturnNewSubcategory{
		Id:           int64(categoryId),
		CategoryName: newSubcategory.CategoryName,
		CategoryId:   int64(newSubcategory.CategoryId),
	}

	resp, err := client.ChangeSubcategory(ctx, grpcNewSubcategory)

	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
