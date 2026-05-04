package recovery

import (
	"authorization-microservice/helpers"
	"authorization-microservice/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/AndreyLebedev1998/auth-grpc"
	"github.com/redis/go-redis/v9"
)

func CodeMatching(w http.ResponseWriter, r *http.Request, db *sql.DB, rdb *redis.Client) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method now allowed", http.StatusMethodNotAllowed)
		return
	}

	var dateRecovery models.DataRecoveryPassword
	ctx := r.Context()

	if err := json.NewDecoder(r.Body).Decode(&dateRecovery); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var grpcParam = &auth.ParamUser{
		Email: dateRecovery.Email,
		Phone: dateRecovery.Phone,
		Id:    0,
	}

	getParamForQuery := helpers.GetParamUser(grpcParam)

	ip := helpers.GetClientIP(r)

	// ключ содержит ip адрес и юзера(контактная информация), чтобы два пользователя с одного ip могли отпралять запрос независимо друг от друга
	key := fmt.Sprintf("recovery_attempts:%s:%s", getParamForQuery, ip)

	attempts, err := rdb.Get(ctx, key).Int()

	if err != nil && err != redis.Nil {
		http.Error(w, "Redis error", http.StatusInternalServerError)
		return
	}

	if attempts >= 3 {
		http.Error(w, "Too many attempts. Try again later.", http.StatusTooManyRequests) // 10 minutes
		return
	}

	query := "SELECT id, code_recovery FROM users WHERE code_recovery_expires_at > NOW() "

	if dateRecovery.Phone != "" || dateRecovery.Email != "" {
		if dateRecovery.Email != "" {
			query += "AND email = $1"
		} else {
			query += "AND phone = $1"
		}

		code := sql.NullString{}
		var userId int

		err := db.QueryRowContext(ctx, query, getParamForQuery).Scan(&userId, &code)

		if err == sql.ErrNoRows {
			http.Error(w, "the code has expired or the user was not found", http.StatusBadRequest)
			return
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if !code.Valid {
			http.Error(w, "code is null", http.StatusBadRequest)
			return
		}

		var codeMatching = code.String == dateRecovery.Code

		if codeMatching {
			query := "UPDATE users SET token_recovery_password = $1, token_recovery_password_expires_at = NOW() + INTERVAL '10 minutes' WHERE id = $2"
			token := helpers.RandomString()

			_, err := db.ExecContext(ctx, query, token, userId)
			if err != nil {
				http.Error(w, "Server error", http.StatusInternalServerError)
				return
			}

			var resp = map[string]string{
				"token": token,
			}

			rdb.Del(ctx, key)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
		} else {
			newAttempts, err := rdb.Incr(ctx, key).Result()
			if err != nil {
				http.Error(w, "Redis error", http.StatusInternalServerError)
				return
			}

			if newAttempts == 1 {
				rdb.Expire(ctx, key, 10*time.Minute)
			}

			time.Sleep(500 * time.Millisecond) // защита от брутфорс атак

			http.Error(w, "the code is not correct", http.StatusBadRequest)
			return
		}
	} else {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
}
