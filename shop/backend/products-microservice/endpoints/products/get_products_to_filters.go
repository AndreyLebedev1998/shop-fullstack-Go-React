package products

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"products-microservice/models"
	"strings"
)

func GetProductsForFilters(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	var minPrice = r.URL.Query().Get("min_price")
	var maxPrice = r.URL.Query().Get("max_price")
	var subcategoryId = r.URL.Query().Get("subcategory_id")
	var categoryId = r.URL.Query().Get("category_id")
	var productsSubcategories []models.ProductsSubcategories

	query := "SELECT products.id, product_name, products.category_id, price, image_url, availability_of_pieces, subcategory_id FROM products JOIN subcategories ON products.subcategory_id = subcategories.id WHERE products.category_id = $1"

	if categoryId == "" {
		fmt.Println(categoryId)
		http.Error(w, "category_id or subcategory_id not defined", http.StatusBadRequest)
		return
	}

	var (
		conditions []string
		args       []interface{}
	)

	args = append(args, categoryId)

	if subcategoryId != "" {
		args = append(args, subcategoryId)
		conditions = append(conditions,
			fmt.Sprintf("products.subcategory_id = $%d", len(args)))
	}

	if minPrice != "" {
		args = append(args, minPrice)
		conditions = append(conditions,
			fmt.Sprintf("price >= $%d", len(args)))
	}

	if maxPrice != "" {
		args = append(args, maxPrice)
		conditions = append(conditions,
			fmt.Sprintf("price <= $%d", len(args)))
	}

	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	for rows.Next() {
		var productsSubcategory models.ProductsSubcategories

		err := rows.Scan(&productsSubcategory.Id, &productsSubcategory.ProductName, &productsSubcategory.CategoryId, &productsSubcategory.Price, &productsSubcategory.ImageUrl,
			&productsSubcategory.AvailabilityOfPieces, &productsSubcategory.SubcategoryId)

		if err != nil {
			http.Error(w, "Error scan row", http.StatusInternalServerError)
			return
		}

		productsSubcategories = append(productsSubcategories, productsSubcategory)
	}

	if productsSubcategories == nil {
		productsSubcategories = make([]models.ProductsSubcategories, 0)
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(productsSubcategories)
}
