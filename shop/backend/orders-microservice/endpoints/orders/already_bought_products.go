package orders

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"orders-microservice/models"
	"slices"

	"github.com/AndreyLebedev1998/auth-grpc"
	productpb "github.com/AndreyLebedev1998/shop-gRPC-product"
)

func AlreadyBoughtProducts(w http.ResponseWriter, r *http.Request, db *sql.DB, clientAuth auth.AuthServiceClient, clientProducts productpb.ProductsServiceClient) {
	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	ctx := r.Context()

	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	respUser, err := clientAuth.GetUserFromToken(ctx, &auth.Token{
		Token: token,
	})

	if err != nil {
		http.Error(w, "Error get user", http.StatusInternalServerError)
		return
	}

	queryGetOrders := `SELECT product_id FROM orders JOIN order_items ON orders.id = order_items.order_id WHERE user_id = $1`

	rows, err := db.QueryContext(ctx, queryGetOrders, respUser.UserId)

	if err == sql.ErrNoRows {
		var empty []interface{}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(empty)
		return
	}

	if err != nil {
		http.Error(w, "Error get products ids", http.StatusInternalServerError)
		return
	}

	var productsIds []int64

	for rows.Next() {
		var productId int64
		err := rows.Scan(&productId)
		if err != nil {
			http.Error(w, "Error scaning product_id", http.StatusInternalServerError)
			return
		}

		productsIds = append(productsIds, productId)
	}

	var productsIdsUniq []int64 = []int64{}

	for _, el := range productsIds {
		if slices.Contains(productsIdsUniq, el) {
			continue
		} else {
			productsIdsUniq = append(productsIdsUniq, el)
		}
	}

	respProducts, err := clientProducts.GetProductsByIds(ctx, &productpb.GetProductsRequest{
		ProductIds: productsIdsUniq,
	})

	if err != nil {
		http.Error(w, "Error get products", http.StatusInternalServerError)
		return
	}

	products := make([]models.ProductResponse, 0, len(respProducts.Products))
	for _, p := range respProducts.Products {
		products = append(products, models.ProductResponse{
			Id:                   p.Id,
			ProductName:          p.ProductName,
			CategoryId:           p.CategoryId,
			CategoryName:         p.CategoryName,
			ImageUrl:             p.ImageUrl,
			AvailabilityOfPieces: p.AvailabilityOfPieces,
			Price:                p.Price,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}
