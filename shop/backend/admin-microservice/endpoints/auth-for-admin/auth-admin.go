package authAdmin

import (
	"admin-microservice/constants"
	"admin-microservice/models"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func Login(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var creds models.AuthAdmin
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Bad request", http.StatusBadGateway)
		return
	}

	var storedHash string
	var adminID int
	ctx := r.Context()

	query := "SELECT id, password FROM admins WHERE nick = $1"

	err := db.QueryRowContext(ctx, query, creds.Nick).Scan(&adminID, &storedHash)

	if err == sql.ErrNoRows {
		http.Error(w, "Invalid name or password", http.StatusUnauthorized)
		return
	}

	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(creds.Password)); err != nil {
		http.Error(w, "Invalid name or password", http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &models.Claims{
		UserID: adminID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(constants.JwtKey)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenString,
		"name":  creds.Nick,
		"id":    strconv.Itoa(adminID),
	})
}
