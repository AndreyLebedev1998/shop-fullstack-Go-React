package products

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"products-microservice/models"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

func GetAllProductsByCategoryId(w http.ResponseWriter, r *http.Request, db *sql.DB, rdb *redis.Client) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var category_id string
	var ctx = r.Context()
	var products []models.Product
	var idStr string = r.URL.Query().Get("category_id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 0 {
		http.Error(w, "Invalid category_id", http.StatusBadRequest)
		return
	}

	if idStr == "" {
		http.Error(w, "category_id cannot be empty", http.StatusBadRequest)
		return
	}

	category_id = idStr

	cacheKey := "products_by_category_id:" + idStr + "all"

	val, err := rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		if json.Unmarshal([]byte(val), &products) == nil {
			w.Header().Add("Content-Type", "application/json")
			json.NewEncoder(w).Encode(products)
			return
		}
	}

	var query string = "SELECT products.id, product_name, products.category_id, price, image_url, availability_of_pieces, subcategory_id FROM products JOIN subcategories ON products.subcategory_id = subcategories.id  WHERE products.category_id = $1"

	rows, err := db.QueryContext(ctx, query, category_id)

	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	for rows.Next() {
		var product models.Product

		err := rows.Scan(&product.Id, &product.ProductName, &product.CategoryId, &product.Price, &product.ImageUrl, &product.AvailabilityOfPieces, &product.SubcategoryId)

		if err != nil {
			fmt.Println(err)
			http.Error(w, "Error while querying the database", http.StatusInternalServerError)
			return
		}

		products = append(products, product)
	}

	w.Header().Set("Content-Type", "application/json")

	bytes, _ := json.Marshal(products)
	rdb.Set(ctx, cacheKey, bytes, 5*time.Minute)

	json.NewEncoder(w).Encode(products)
}

func GetAllProductsHandler(db *sql.DB, rdb *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		GetAllProductsByCategoryId(w, r, db, rdb)
	}
}

// @Summary find products by category_id
// @Description Returns products by category_id
// @Tags products
// @Accept json
// @Produce json
// @Param id query int true "category_id"
// @Success 200 {array} models.Product
// @Failure 400 {string} string
// @Failure 405 {string} string "Method not allowed"
// @Router /products [get]
func GetAllProductsForSwagger(w http.ResponseWriter, r *http.Request) {
	// Пустое тело, только для Swagger
}
