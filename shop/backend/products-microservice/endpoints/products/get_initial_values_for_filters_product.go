package products

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"products-microservice/models"
)

func GetInitialValuesForFilter(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	var initialValues models.InitialValuesForFilter
	var categoryId = r.URL.Query().Get("category_id")
	var subcategoryId = r.URL.Query().Get("subcategory_id")
	if categoryId == "" {
		fmt.Println(categoryId)
		http.Error(w, "category_id not defined", http.StatusBadRequest)
		return
	}
	query := "SELECT COALESCE(MIN(price), 0), COALESCE(MAX(price), 0) FROM products WHERE category_id = $1"

	var args []any
	args = append(args, categoryId)
	if subcategoryId != "" {
		query += " AND subcategory_id = $2"
		args = append(args, subcategoryId)
	}

	err := db.QueryRowContext(ctx, query, args...).Scan(&initialValues.MinPrice, &initialValues.MaxPrice)

	if err == sql.ErrNoRows {
		initialValues.MinPrice = 0
		initialValues.MaxPrice = 0

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(initialValues)
		return
	}

	if err != nil {
		fmt.Println(err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(initialValues)
}
