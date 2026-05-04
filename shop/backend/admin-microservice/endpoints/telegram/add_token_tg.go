package telegram

import (
	"admin-microservice/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func AddTokenTg(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var token models.NewTelegramToken
	ctx := r.Context()

	if err := json.NewDecoder(r.Body).Decode(&token); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	query := "INSERT INTO telegram_tokens (id, token) VALUES (1, $1) ON CONFLICT (id) DO UPDATE SET token = EXCLUDED.token RETURNING id"

	var id int

	err := db.QueryRowContext(ctx, query, token.Token).Scan(&id)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"id":    strconv.Itoa(id),
		"token": token.Token,
	})
}
