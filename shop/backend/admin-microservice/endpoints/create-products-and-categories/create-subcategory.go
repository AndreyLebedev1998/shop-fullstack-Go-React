package create

import (
	"admin-microservice/models"
	"encoding/json"
	"fmt"
	"net/http"

	product "github.com/AndreyLebedev1998/shop-gRPC-product"
)

func CreateSubcategory(w http.ResponseWriter, r *http.Request, client product.ProductsServiceClient) {
	if r.Method != http.MethodPost {
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

	var grpcNewCategory = &product.NewSubcategory{
		CategoryName: newSubcategory.CategoryName,
		CategoryId:   int64(newSubcategory.CategoryId),
	}

	resp, err := client.CreateSubcategory(ctx, grpcNewCategory)

	if err != nil {
		fmt.Println(err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
