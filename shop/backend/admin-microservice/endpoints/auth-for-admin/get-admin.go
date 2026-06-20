package authAdmin

import (
	"admin-microservice/constants"
	"admin-microservice/models"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
)

func GetAdmin(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Извлекаем токен
	tokenStr := r.Header.Get("Authorization")
	if tokenStr == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	ctx := r.Context()
	if len(tokenStr) > 7 && tokenStr[:7] == "Bearer " {
		tokenStr = tokenStr[7:]
	}

	claims := &models.Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return constants.JwtKey, nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var nick string

	err = db.QueryRowContext(ctx, "SELECT nick FROM admins WHERE id = $1", claims.UserID).Scan(&nick)

	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]any{
		"authenticated": true,
		"user_id":       claims.UserID,
		"nick":          nick,
	})
}
