package create

import (
	"admin-microservice/models"
	"database/sql"
	"encoding/json"
	"net/http"
)

func CreateCategory(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var categoryName models.NewCategory
	var category models.Category
	var ctx = r.Context()

	if err := json.NewDecoder(r.Body).Decode(&categoryName); err != nil {
		http.Error(w, "category_name is not defined", http.StatusBadRequest)
		return
	}

	if categoryName.CategoryName == "" {
		http.Error(w, "category_name can't be empty", http.StatusBadRequest)
		return
	}

	query := "INSERT INTO categories (category_name) VALUES ($1) RETURNING id"

	err := db.QueryRowContext(ctx, query, categoryName.CategoryName).Scan(&category.Id)
	if err != nil {
		http.Error(w, "Error inserting into the database", http.StatusInternalServerError)
		return
	}

	category.CategoryName = categoryName.CategoryName

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}
