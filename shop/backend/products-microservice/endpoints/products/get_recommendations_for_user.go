package products

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"products-microservice/models"
	"time"

	"github.com/AndreyLebedev1998/auth-grpc"
	order "github.com/AndreyLebedev1998/shop-gRPC-orders"
	"github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

func GetRecommendationsForUser(w http.ResponseWriter, r *http.Request, db *sql.DB, rdb *redis.Client, clientAuth auth.AuthServiceClient, clientOrders order.OrderServiceClient) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	ctx := r.Context()

	respAuth, err := clientAuth.GetUserFromToken(ctx, &auth.Token{
		Token: token,
	})

	if err != nil {
		http.Error(w, "Error get user from token", http.StatusUnauthorized)
		return
	}

	cacheKey := fmt.Sprintf("recommendations_for_user:%d", respAuth.UserId)
	var products []models.RecommendationProduct = make([]models.RecommendationProduct, 0)

	val, err := rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		if json.Unmarshal([]byte(val), &products) == nil {
			w.Header().Add("Content-Type", "application/json")
			json.NewEncoder(w).Encode(products)
			return
		}
	}

	respOrder, err := clientOrders.GetProductsIdsFromOrdersForUsers(ctx, &order.UserId{
		UserId: respAuth.UserId,
	})

	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error get user product orders", http.StatusInternalServerError)
		return
	}

	var productsIds []int64
	seen := make(map[int64]struct{})

	for _, p := range respOrder.ProductsIds {
		if _, ok := seen[p.Id]; !ok {
			seen[p.Id] = struct{}{}
			productsIds = append(productsIds, p.Id)
		}
	}

	query := "SELECT DISTINCT category_id FROM products WHERE id = ANY($1)"

	rows, err := db.QueryContext(ctx, query, pq.Array(productsIds))

	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error get categories for product ids user", http.StatusInternalServerError)
		return
	}

	var categoryIds []int
	seenCategoryId := make(map[int]struct{})

	for rows.Next() {
		var categoryId int

		err := rows.Scan(&categoryId)
		if err != nil {
			http.Error(w, "Error scan category_id", http.StatusInternalServerError)
			return
		}

		if _, ok := seenCategoryId[categoryId]; !ok {
			seenCategoryId[categoryId] = struct{}{}
			categoryIds = append(categoryIds, categoryId)
		}
	}

	queryRecommendations := `SELECT * FROM (
    							SELECT
									products.id,
        							product_name,
									price,
        							category_id, 
									category_name,
									image_url,
									availability_of_pieces,
									subcategory_id,
        							ROW_NUMBER() OVER (PARTITION BY category_id ORDER BY purchase_count DESC) AS rn
    								FROM products 
    								JOIN product_stats ON products.id = product_stats.product_id
									JOIN categories ON categories.id = products.category_id
									) sub
							WHERE rn <= 5 AND category_id = ANY($1)`

	rows, err = db.QueryContext(ctx, queryRecommendations, pq.Array(categoryIds))
	if err == sql.ErrNoRows {
		empty := make([]any, 0, 2)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(empty)
	}
	if err != nil {
		http.Error(w, "Error get recommendations", http.StatusInternalServerError)
		return
	}

	for rows.Next() {
		var product models.RecommendationProduct
		err := rows.Scan(&product.Id, &product.ProductName, &product.Price, &product.CategoryId, &product.CategoryName, &product.ImageUrl, &product.AvailabilityOfPieces, &product.SubcategoryId, &product.Rating)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Error scanning recommendation product", http.StatusInternalServerError)
			return
		}

		if product.ImageUrl == nil {
			product.ImageUrl = new(string)
		}
		products = append(products, product)
	}

	if len(products) > 0 {
		bytes, _ := json.Marshal(products)
		rdb.Set(ctx, cacheKey, bytes, 10*time.Minute)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}
