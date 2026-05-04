package recovery

import (
	"authorization-microservice/helpers"
	"authorization-microservice/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AndreyLebedev1998/auth-grpc"
	"golang.org/x/crypto/bcrypt"
)

func RecoveryPassword(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var recovery models.RecoveryPassword
	ctx := r.Context()

	if err := json.NewDecoder(r.Body).Decode(&recovery); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if recovery.NewPassword != recovery.ConfirmationPassword {
		http.Error(w, "passwords do not match", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(recovery.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	query := "UPDATE users SET password = $1 WHERE token_recovery_password_expires_at > NOW() AND token_recovery_password = $2 AND "

	var grpcParam = &auth.ParamUser{
		Email: recovery.Email,
		Phone: recovery.Phone,
		Id:    0,
	}

	getParamForQuery := helpers.GetParamUser(grpcParam)

	if recovery.Email != "" || recovery.Phone != "" {
		if recovery.Email != "" {
			query += "email = $3"
		} else {
			query += "phone = $3"
		}

		rows, err := db.ExecContext(ctx, query, hashedPassword, recovery.Token, getParamForQuery)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		rowsUpdated, _ := rows.RowsAffected()

		if rowsUpdated != 1 {
			http.Error(w, "The user was not found or token has expired", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Password changed successfully",
		})
	} else {
		http.Error(w, "Email or phone is not defined", http.StatusBadRequest)
		return
	}
}
