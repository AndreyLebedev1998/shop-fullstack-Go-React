package products

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/AndreyLebedev1998/auth-grpc"
)

func RemoveFavoriteProduct(w http.ResponseWriter, r *http.Request, db *sql.DB, clientAuth auth.AuthServiceClient) {
	if r.Method != http.MethodDelete {
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

	var productIdStr = r.URL.Query().Get("product_id")
	productId, err := strconv.Atoi(productIdStr)
	if err != nil || productId <= 0 {
		http.Error(w, "Invalid product_id", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	respAuth, err := clientAuth.GetUserFromToken(ctx, &auth.Token{
		Token: token,
	})

	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	query := "DELETE FROM favorite_products WHERE user_id = $1 AND product_id = $2"

	resp, err := db.ExecContext(ctx, query, respAuth.UserId, productId)
	if err != nil {
		http.Error(w, "Error remove favorite product", http.StatusInternalServerError)
		return
	}

	rowsDeleted, _ := resp.RowsAffected()

	if rowsDeleted != 1 {
		http.Error(w, "favorite product is not defined", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"user_id":    strconv.FormatInt(respAuth.UserId, 10),
		"product_id": productIdStr,
	})
}
