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
	"time"

	"github.com/redis/go-redis/v9"
)

func GetAllProductsBySubcategoryId(w http.ResponseWriter, r *http.Request, db *sql.DB, rdb *redis.Client) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var subcategoryIdStr = r.URL.Query().Get("subcategory_id")
	var productsSubcategories []models.ProductsSubcategories
	ctx := r.Context()
	subcategoryId, err := strconv.Atoi(subcategoryIdStr)
	if err != nil || subcategoryId <= 0 {
		http.Error(w, "Invalid subcategory_id", http.StatusBadRequest)
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

	var cacheKey string

	if indicator != "" {
		cacheKey = fmt.Sprintf("products:by-subcategory::id:%d:indicator:%s", subcategoryId, indicator)
	} else {
		cacheKey = fmt.Sprintf("products:by-subcategory:id:%d", subcategoryId)
	}

	val, err := rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		if json.Unmarshal([]byte(val), &productsSubcategories) == nil {
			fmt.Println("Redis")
			w.Header().Add("Content-Type", "application/json")
			json.NewEncoder(w).Encode(productsSubcategories)
			return
		}
	}

	query := `SELECT products.id, product_name, products.category_id, price, image_url, availability_of_pieces, subcategory_id, category_name FROM products 
			  JOIN subcategories ON subcategories.id = products.subcategory_id WHERE products.subcategory_id = $1 `

	if indicator != "" {
		query += helpers.GetSortForIndicator(constants.Indicator(indicator))
	} else {
		query += "ORDER BY availability_of_pieces DESC"
	}

	rows, err := db.QueryContext(ctx, query, subcategoryId)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	for rows.Next() {
		var productsSubcategory models.ProductsSubcategories

		err := rows.Scan(&productsSubcategory.Id, &productsSubcategory.ProductName, &productsSubcategory.CategoryId, &productsSubcategory.Price, &productsSubcategory.ImageUrl,
			&productsSubcategory.AvailabilityOfPieces, &productsSubcategory.SubcategoryId, &productsSubcategory.SubcategoryName)

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

	bytes, _ := json.Marshal(productsSubcategories)
	rdb.Set(ctx, cacheKey, bytes, 5*time.Minute)

	json.NewEncoder(w).Encode(productsSubcategories)
}
