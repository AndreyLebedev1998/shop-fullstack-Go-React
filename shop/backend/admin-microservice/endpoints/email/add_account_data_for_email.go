package email

import (
	"admin-microservice/models"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
)

func AddAccountDataForEmail(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var emailConnData models.EmailConnData
	ctx := r.Context()

	if err := json.NewDecoder(r.Body).Decode(&emailConnData); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if emailConnData.SenderEmail == "" || emailConnData.Password == "" {
		http.Error(w, "email or password is not defined", http.StatusBadRequest)
		return
	}

	query := "INSERT INTO email_conn (id, sender_email, password) VALUES (1, $1, $2) ON CONFLICT (id) DO UPDATE SET sender_email = EXCLUDED.sender_email, password = EXCLUDED.password"

	rows, err := db.ExecContext(ctx, query, emailConnData.SenderEmail, emailConnData.Password)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	rowsUpdated, _ := rows.RowsAffected()

	if rowsUpdated != 1 {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	var result = map[string]string{
		"id":           strconv.Itoa(1),
		"sender_email": emailConnData.SenderEmail,
		"message":      "Mail connection data has been updated successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
