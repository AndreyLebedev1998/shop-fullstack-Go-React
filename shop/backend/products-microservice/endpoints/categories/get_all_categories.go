package categories

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"products-microservice/models"

	"github.com/redis/go-redis/v9"
)

func GetAllCategories(w http.ResponseWriter, r *http.Request, db *sql.DB, rdb *redis.Client) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var categories []models.CategoryWithSubcategories
	var categoriesMap = make(map[int]*models.CategoryWithSubcategories)
	var ctx = r.Context()

	/* cacheKey := "categories:all"

	val, err := rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		if json.Unmarshal([]byte(val), &categoriesMap) == nil {
			w.Header().Add("Content-Type", "application/json")
			fmt.Println("Redis")
			json.NewEncoder(w).Encode(categoriesMap)
			return
		}
	} */

	var query string = `SELECT categories.id, categories.category_name, subcategories.category_name AS subcategory_name, subcategories.id AS subcategory_id
						FROM categories JOIN subcategories ON categories.id = subcategories.category_id`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		http.Error(w, "Error while querying the database", http.StatusInternalServerError)
		return
	}

	for rows.Next() {
		var category models.CategoryWithSubcategories
		var subcategorie models.Subcategory
		err := rows.Scan(&category.Id, &category.CategoryName, &subcategorie.CategoryName, &subcategorie.SubcategoryId)
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		subcategorie.CategoryId = category.Id

		if _, inMap := categoriesMap[category.Id]; !inMap {
			categoriesMap[category.Id] = &category
			categoriesMap[category.Id].Subcategories = append(categoriesMap[category.Id].Subcategories, subcategorie)
		} else {
			categoriesMap[category.Id].Subcategories = append(categoriesMap[category.Id].Subcategories, subcategorie)
		}
	}

	for _, value := range categoriesMap {
		categories = append(categories, *value)
	}

	/* bytes, _ := json.Marshal(categories)
	rdb.Set(ctx, cacheKey, bytes, 5*time.Minute) */

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}
