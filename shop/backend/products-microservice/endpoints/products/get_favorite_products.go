package products

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"products-microservice/models"

	"github.com/AndreyLebedev1998/auth-grpc"
)

func GetFavoriteProducts(w http.ResponseWriter, r *http.Request, db *sql.DB, clientAuth auth.AuthServiceClient) {
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

	var favoriteProducts []models.AllFavoriteProduct

	query := `SELECT products.id, user_id, product_name, price, category_id, image_url, availability_of_pieces, subcategory_id FROM favorite_products
			  JOIN products ON favorite_products.product_id = products.id
			  WHERE user_id = $1 ORDER BY availability_of_pieces DESC`

	rows, err := db.QueryContext(ctx, query, respAuth.UserId)

	if err == sql.ErrNoRows {
		var empty = make([]interface{}, 0)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(empty)
	}

	if err != nil {
		http.Error(w, "Error geting favorite products", http.StatusUnauthorized)
		return
	}

	for rows.Next() {
		var favoriteProduct models.AllFavoriteProduct

		err := rows.Scan(&favoriteProduct.Id, &favoriteProduct.UserId, &favoriteProduct.ProductName, &favoriteProduct.Price,
			&favoriteProduct.CategoryId, &favoriteProduct.ImageUrl, &favoriteProduct.AvailabilityOfPieces, &favoriteProduct.SubcategoryId)

		if err != nil {
			http.Error(w, "Error scaning favorite product", http.StatusInternalServerError)
			return
		}

		favoriteProducts = append(favoriteProducts, favoriteProduct)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(favoriteProducts)
}
