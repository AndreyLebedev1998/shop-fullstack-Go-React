package change

import (
	"authorization-microservice/helpers"
	"authorization-microservice/models"
	"database/sql"
	"encoding/json"
	"net/http"
)

func ChangeAccount(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var user models.UserEntrance
	ctx := r.Context()

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	userId, err := helpers.GetUserIDFromToken(r)

	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	query := "UPDATE users SET name = $1, lastname = $2, email = $3, phone = $4 WHERE id = $5"

	rows, err := db.ExecContext(ctx, query, user.Name, user.LastName, user.Email, user.Phone, userId)
	if err != nil {
		http.Error(w, "Server error", http.StatusMethodNotAllowed)
		return
	}

	rowsUpdated, _ := rows.RowsAffected()

	if rowsUpdated != 1 {
		http.Error(w, "user is not defined", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "user data changed successfully",
	})
}
