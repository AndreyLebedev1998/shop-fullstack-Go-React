package autentification

import (
	"authorization-microservice/constants"
	"authorization-microservice/models"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Login(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var creds models.Entrance
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Bad request", http.StatusBadGateway)
		return
	}

	var storedHash string
	var userID int
	err := db.QueryRow("SELECT id, password FROM users WHERE name=$1", creds.Name).Scan(&userID, &storedHash)

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
		UserID: userID,
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

	var user models.UserEntrance

	err = db.QueryRow("SELECT name, lastname, email, phone FROM users WHERE name = $1", creds.Name).Scan(&user.Name, &user.LastName, &user.Email, &user.Phone)

	if err != nil {
		http.Error(w, "Server error. Error get user", http.StatusInternalServerError)
		return
	}

	var userEmail string
	if user.Email != nil {
		userEmail = *user.Email
	}

	var userPhone string

	if user.Phone != nil {
		userPhone = *user.Phone
	}

	json.NewEncoder(w).Encode(map[string]string{
		"token":    tokenString,
		"name":     user.Name,
		"lastname": user.LastName,
		"email":    userEmail,
		"phone":    userPhone,
	})
}

func Me(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Извлекаем токен
	tokenStr := r.Header.Get("Authorization")
	if tokenStr == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
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

	var user string

	err = db.QueryRow("SELECT name FROM users WHERE id = $1", claims.UserID).Scan(&user)

	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	var userEntrance models.UserEntrance

	err = db.QueryRow("SELECT name, lastname, email, phone FROM users WHERE name = $1", user).Scan(&userEntrance.Name, &userEntrance.LastName, &userEntrance.Email, &userEntrance.Phone)

	if err != nil {
		http.Error(w, "Server error. Error get user", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]any{
		"authenticated": true,
		"user_id":       claims.UserID,
		"name":          userEntrance.Name,
		"lastname":      userEntrance.LastName,
		"email":         userEntrance.Email,
		"phone":         userEntrance.Phone,
	})
}
