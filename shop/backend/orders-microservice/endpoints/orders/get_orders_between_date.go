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

func GetOrdersByUserBetweenDate(w http.ResponseWriter, r *http.Request, db *sql.DB, rdb *redis.Client, client product.ProductsServiceClient) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var emailParam = r.URL.Query().Get("email")
	var phoneParam = r.URL.Query().Get("phone")
	var userIdParam = r.URL.Query().Get("user_id")
	var startDateParam = r.URL.Query().Get("start_date")
	var endDateParam = r.URL.Query().Get("end_date")
	var productIDsSet = make(map[int]struct{})
	var orders []models.FullOrder
	var ordersMap = make(map[int]*models.FullOrder)
	var ctx = r.Context()
	if startDateParam == "" {
		http.Error(w, "date can't be empty", http.StatusBadRequest)
		return
	}
	_, err := time.Parse("2006-01-02", startDateParam)
	if err != nil {
		http.Error(w, "invalid date format, use YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	cacheKey := "orders:between-date:user"

	val, err := rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		if json.Unmarshal([]byte(val), &orders) == nil {
			w.Header().Add("Content-Type", "application/json")
			fmt.Println("Redis")
			json.NewEncoder(w).Encode(orders)
			return
		}
	}

	var query string = `SELECT orders.id as order_id, user_id, email, phone, status, total_price, created_at, 
						product_id, quantity, order_items.price
						FROM orders 
						JOIN order_items ON orders.id = order_items.order_id `
	if emailParam != "" || phoneParam != "" || userIdParam != "" {
		var params = []models.ParamsForQuery{
			{
				Column: "email",
				Value:  emailParam,
			},
			{
				Column: "phone",
				Value:  phoneParam,
			},
			{
				Column: "user_id",
				Value:  userIdParam,
			},
		}

		dynamicQuery, args := helpers.SqlQueryWithParamAndDate(params)
		allArgs := append([]any{startDateParam, endDateParam}, args...)
		rows, err := db.QueryContext(ctx, query+" "+dynamicQuery, allArgs...)

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
			orders = append(orders, *order)
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

	for i := range orders {
		for j := range orders[i].Products {

			pid := int64(orders[i].Products[j].ProductId)

			if p, ok := productsMap[pid]; ok {
				orders[i].Products[j] = models.Products{
					ProductId:    int(p.Id),
					ProductName:  p.ProductName,
					CategoryId:   int(p.CategoryId),
					CategoryName: p.CategoryName,
					ImageUrl:     &p.ImageUrl,
				}
			}
		}
	}

	bytes, _ := json.Marshal(orders)
	rdb.Set(ctx, cacheKey, bytes, 5*time.Minute)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}
