package orders

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"orders-microservice/helpers"
	"orders-microservice/models"
	"time"

	"github.com/AndreyLebedev1998/auth-grpc"
	product "github.com/AndreyLebedev1998/shop-gRPC-product"
	"github.com/redis/go-redis/v9"
)

func GetOrdersByParametr(w http.ResponseWriter, r *http.Request, db *sql.DB, rdb *redis.Client, client product.ProductsServiceClient, clientAuth auth.AuthServiceClient) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var token models.Token

	if err := json.NewDecoder(r.Body).Decode(&token); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	var ctx = r.Context()

	respUser, err := clientAuth.GetUserFromToken(ctx, &auth.Token{
		Token: token.Token,
	})

	if err != nil {
		http.Error(w, "Error get user id", http.StatusInternalServerError)
		return
	}

	var userId int64

	if respUser != nil {
		userId = respUser.UserId
	} else {
		http.Error(w, "Error get user id", http.StatusInternalServerError)
		return
	}

	var fullOrders []models.FullOrder
	ordersMap := make(map[int]*models.FullOrder)
	var productIDsSet = make(map[int]struct{})

	cacheKey := fmt.Sprintf("orders:orders:user_id:%d", userId)

	val, err := rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		if json.Unmarshal([]byte(val), &fullOrders) == nil {
			w.Header().Set("Content-Type", "application/json")
			fmt.Println("Redis")
			json.NewEncoder(w).Encode(fullOrders)
			return
		}
	}

	var query string = `SELECT orders.id as order_id, user_id, email, phone, status, total_price, created_at, updated_at,
						product_id, quantity, order_items.price
						FROM orders 
						JOIN order_items ON orders.id = order_items.order_id 
						WHERE user_id = $1`

	rows, err := db.QueryContext(ctx, query, userId)

	if err != nil {
		http.Error(w, "Error while querying the database", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	if err := helpers.ForRowsAfterQuery(rows, ordersMap, productIDsSet); err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	for _, order := range ordersMap {
		fullOrders = append(fullOrders, *order)
	}

	var productIDs []int64

	for id := range productIDsSet {
		productIDs = append(productIDs, int64(id))
	}

	resp, err := client.GetProductsByIds(ctx, &product.GetProductsRequest{
		ProductIds: productIDs,
	})

	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	productsMap := make(map[int64]*product.Product)

	for _, p := range resp.Products {
		productsMap[p.Id] = p
	}

	for i := range fullOrders {
		fmt.Println(fullOrders)
		for j := range fullOrders[i].Products {

			pid := int64(fullOrders[i].Products[j].ProductId)

			if p, ok := productsMap[pid]; ok {
				fullOrders[i].Products[j].ProductName = p.ProductName
				fullOrders[i].Products[j].CategoryId = int(p.CategoryId)
				fullOrders[i].Products[j].CategoryName = p.CategoryName
				fullOrders[i].Products[j].ImageUrl = &p.ImageUrl
			}
		}
	}

	if fullOrders == nil {
		fullOrders = []models.FullOrder{}
	}

	bytes, _ := json.Marshal(fullOrders)
	rdb.Set(ctx, cacheKey, bytes, 5*time.Minute)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fullOrders)
}
