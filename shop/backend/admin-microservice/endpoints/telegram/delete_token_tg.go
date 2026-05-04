package telegram

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

func RemoveToken(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	query := "DELETE FROM telegram_tokens WHERE id = 1"

	resp, err := db.ExecContext(ctx, query)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	rowsDeleted, _ := resp.RowsAffected()

	if rowsDeleted != 1 {
		http.Error(w, "Token is not defined", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Token deleted successfully",
	})
}
