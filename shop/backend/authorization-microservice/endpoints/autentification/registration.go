package autentification

import (
	"authorization-microservice/models"
	"database/sql"
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func Register(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var creds models.Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	email := sql.NullString{}
	phone := sql.NullString{}
	if creds.Email != nil && *creds.Email != "" {
		email = sql.NullString{
			String: *creds.Email,
			Valid:  true,
		}
	}

	if creds.Phone != nil && *creds.Phone != "" {
		phone = sql.NullString{
			String: *creds.Phone,
			Valid:  true,
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	var id int
	err = db.QueryRow(`INSERT INTO users (name, lastname, password, email, phone)
					  VALUES ($1, $2, $3, $4, $5) RETURNING id
	`, creds.Name, creds.LastName, hashedPassword, email, phone).Scan(&id)

	if err != nil {
		http.Error(w, "email or phone already exists", http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"user_id": id})
}
