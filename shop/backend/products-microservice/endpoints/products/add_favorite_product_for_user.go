package products

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"products-microservice/models"

	"github.com/AndreyLebedev1998/auth-grpc"
)

func AddFavoriteProductForUser(w http.ResponseWriter, r *http.Request, db *sql.DB, clientAuth auth.AuthServiceClient) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

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
		http.Error(w, "Error get user", http.StatusUnauthorized)
		return
	}

	var productId models.ProductId

	if err := json.NewDecoder(r.Body).Decode(&productId); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	queryInsert := "INSERT INTO favorite_products (user_id, product_id) VALUES ($1, $2)"

	rows, err := db.ExecContext(ctx, queryInsert, respUser.UserId, productId.ProductId)

	if err != nil {
		http.Error(w, "Error inserting data", http.StatusInternalServerError)
		return
	}

	rowsInsert, _ := rows.RowsAffected()

	if rowsInsert != 1 {
		http.Error(w, "Data not saved", http.StatusInternalServerError)
		return
	}

	var favoriteProduct models.AllFavoriteProduct

	queryGetProduct := "SELECT products.id, product_name, category_id, price, image_url, subcategory_id FROM products WHERE products.id = $1"

	err = db.QueryRowContext(ctx, queryGetProduct, productId.ProductId).Scan(&favoriteProduct.Id, &favoriteProduct.ProductName,
		&favoriteProduct.CategoryId, &favoriteProduct.Price, &favoriteProduct.ImageUrl, &favoriteProduct.SubcategoryId)

	if err != nil {
		http.Error(w, "Error geting product. The product has been saved.", http.StatusInternalServerError)
		return
	}

	if favoriteProduct.ImageUrl == nil {
		empty := ""
		favoriteProduct.ImageUrl = &empty
	}

	favoriteProduct.UserId = int(respUser.UserId)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(favoriteProduct)
}
