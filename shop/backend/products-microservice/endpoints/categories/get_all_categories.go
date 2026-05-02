package categories

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"products-microservice/models"
	"time"

	"github.com/redis/go-redis/v9"
)

func GetAllCategories(w http.ResponseWriter, r *http.Request, db *sql.DB, rdb *redis.Client) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var categories []models.Category
	var ctx = r.Context()

	cacheKey := "categories:all"

	val, err := rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		if json.Unmarshal([]byte(val), &categories) == nil {
			w.Header().Add("Content-Type", "application/json")
			fmt.Println("Redis")
			json.NewEncoder(w).Encode(categories)
			return
		}
	}

	var query string = "SELECT * FROM categories"

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		http.Error(w, "Error while querying the database", http.StatusInternalServerError)
		return
	}

	for rows.Next() {
		var category models.Category
		err := rows.Scan(&category.Id, &category.CategoryName)
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		categories = append(categories, category)
	}

	bytes, _ := json.Marshal(categories)
	rdb.Set(ctx, cacheKey, bytes, 5*time.Minute)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}
