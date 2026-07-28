package products

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"products-microservice/constants"
	"products-microservice/helpers"
	"products-microservice/models"
	"strconv"
	"strings"
)

func SortProductsForIndicator(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var indicator string = r.URL.Query().Get("indicator")
	isValid := helpers.IsValidIndicator(constants.Indicator(indicator))

	if !isValid {
		http.Error(w, "Indicator is not valid", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	query := `SELECT products.id, product_name, products.category_id, price, image_url, availability_of_pieces, subcategory_id, category_name FROM products 
			  JOIN categories ON products.category_id = categories.id `

	var conditions []string
	var queryArgs []any

	var categoryIdStr string = r.URL.Query().Get("category_id")
	if categoryIdStr != "" {
		categoryId, err := strconv.Atoi(categoryIdStr)
		if err != nil {
			http.Error(w, "Invalid category_id", http.StatusBadRequest)
			return
		}
		conditions = append(conditions, fmt.Sprintf("products.category_id = $%d", len(queryArgs)+1))
		queryArgs = append(queryArgs, categoryId)
	}

	var subcategoryIdStr string = r.URL.Query().Get("subcategory_id")
	if subcategoryIdStr != "" {
		subcategoryId, err := strconv.Atoi(subcategoryIdStr)
		if err != nil {
			http.Error(w, "Invalid subcategory_id", http.StatusBadRequest)
			return
		}
		conditions = append(conditions, fmt.Sprintf("subcategory_id = $%d", len(queryArgs)+1))
		queryArgs = append(queryArgs, subcategoryId)
	}

	if len(conditions) > 0 {
		query += "WHERE " + strings.Join(conditions, " AND ") + " "
	}

	query += helpers.GetSortForIndicator(constants.Indicator(indicator))

	rows, err := db.QueryContext(ctx, query, queryArgs...)
	if err == sql.ErrNoRows {
		empty := make([]any, 0)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(empty)
		return
	}

	if err != nil {
		http.Error(w, "Error get products", http.StatusInternalServerError)
		return
	}

	var products []models.Product

	for rows.Next() {
		var product models.Product

		err := rows.Scan(&product.Id, &product.ProductName, &product.CategoryId, &product.Price, &product.ImageUrl, &product.AvailabilityOfPieces, &product.SubcategoryId, &product.CategoryName)

		if err != nil {
			fmt.Println(err)
			http.Error(w, "Error while querying the database", http.StatusInternalServerError)
			return
		}

		products = append(products, product)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}
