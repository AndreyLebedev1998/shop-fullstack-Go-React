package orders

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"orders-microservice/helpers"
	"orders-microservice/models"
	"time"

	product "github.com/AndreyLebedev1998/shop-gRPC-product"
	"github.com/redis/go-redis/v9"
)

func GetOrdersByParametr(w http.ResponseWriter, r *http.Request, db *sql.DB, rdb *redis.Client, client product.ProductsServiceClient) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var emailParametr string = r.URL.Query().Get("email")
	var phoneParamert string = r.URL.Query().Get("phone")
	var userIdParametr string = r.URL.Query().Get("user_id")
	var ctx = r.Context()
	var fullOrders []models.FullOrder
	ordersMap := make(map[int]*models.FullOrder)
	var productIDsSet = make(map[int]struct{})

	cacheKey := "orders:user"

	val, err := rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		if json.Unmarshal([]byte(val), &fullOrders) == nil {
			w.Header().Add("Content-Type", "application/json")
			fmt.Println("Redis")
			json.NewEncoder(w).Encode(fullOrders)
			return
		}
	}

	var query string = `SELECT orders.id as order_id, user_id, email, phone, status, total_price, created_at, 
						product_id, quantity, order_items.price
						FROM orders 
						JOIN order_items ON orders.id = order_items.order_id`

	if emailParametr != "" || phoneParamert != "" || userIdParametr != "" {
		var params = []models.ParamsForQuery{
			{
				Column: "email",
				Value:  emailParametr,
			},
			{
				Column: "phone",
				Value:  phoneParamert,
			},
			{
				Column: "user_id",
				Value:  userIdParametr,
			},
		}

		dynamicQuery, args := helpers.SqlQueryWithParam(params)

		rows, err := db.QueryContext(ctx, query+" "+dynamicQuery, args...)

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
		for j := range fullOrders[i].Products {

			pid := int64(fullOrders[i].Products[j].ProductId)

			if p, ok := productsMap[pid]; ok {
				fullOrders[i].Products[j] = models.Products{
					ProductId:    int(p.Id),
					ProductName:  p.ProductName,
					CategoryId:   int(p.CategoryId),
					CategoryName: p.CategoryName,
					ImageUrl:     &p.ImageUrl,
				}
			}
		}
	}

	bytes, _ := json.Marshal(fullOrders)
	rdb.Set(ctx, cacheKey, bytes, 5*time.Minute)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fullOrders)
}
