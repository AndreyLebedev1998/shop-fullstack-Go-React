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

func FindProductBySymbols(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	symbols := r.URL.Query().Get("symbols")
	var idStr string = r.URL.Query().Get("category_id")
	if idStr != "" {
		if _, err := strconv.Atoi(idStr); err != nil {
			http.Error(w, "Invalid category_id", http.StatusBadRequest)
			return
		}
	}

	if symbols == "" || strings.TrimSpace(symbols) == "" {
		var empty []any = make([]any, 0)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(empty)
		return
	}

	var indicator string = r.URL.Query().Get("indicator")
	isValid := helpers.IsValidIndicator(constants.Indicator(indicator))

	if indicator != "" {
		if !isValid {
			http.Error(w, "Indicator is not valid", http.StatusBadRequest)
			return
		}
	}

	ctx := r.Context()

	query := `SELECT products.id, product_name, products.category_id, price, image_url, availability_of_pieces, subcategory_id, category_name FROM products 
						JOIN categories ON products.category_id = categories.id 
						WHERE product_name ILIKE $1 OR EXISTS (SELECT 1 FROM unnest(tags) AS tag WHERE tag ILIKE $2) `

	allArgs := append([]any{}, "%"+symbols+"%", "%"+symbols+"%")
	if idStr != "" {
		query += " AND products.category_id = $3"
		allArgs = append(allArgs, idStr)
	}

	if indicator != "" {
		query += helpers.GetSortForIndicator(constants.Indicator(indicator))
	} else {
		query += " ORDER BY availability_of_pieces DESC"
	}

	rows, err := db.QueryContext(ctx, query, allArgs...)

	if err != nil {
		http.Error(w, "Error get products", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

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
	if len(products) == 0 {
		var empty []any = make([]any, 0)
		json.NewEncoder(w).Encode(empty)
		return
	}

	json.NewEncoder(w).Encode(products)
}
