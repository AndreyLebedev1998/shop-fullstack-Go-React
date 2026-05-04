package recovery

import (
	"authorization-microservice/helpers"
	"authorization-microservice/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/smtp"
	"time"

	"github.com/AndreyLebedev1998/admin-grpc"
	"github.com/AndreyLebedev1998/auth-grpc"
	"github.com/redis/go-redis/v9"
)

func SendCodeRecoveryPasswordInEmail(w http.ResponseWriter, r *http.Request, db *sql.DB, rdb *redis.Client, clientAdmin admin.AdminServiceClient) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var recoveryPasswordData models.RecoveryPasswordData
	ctx := r.Context()

	if err := json.NewDecoder(r.Body).Decode(&recoveryPasswordData); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	resp, err := clientAdmin.GetDataFromConnectionEmail(ctx, &admin.Empty{})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if recoveryPasswordData.Email == "" {
		http.Error(w, "email is not defined", http.StatusBadRequest)
		return
	}

	var email string
	var userId int
	queryGetChatId := "SELECT id, email FROM users WHERE email = $1"

	err = db.QueryRowContext(ctx, queryGetChatId, recoveryPasswordData.Email).Scan(&userId, &email)
	if err == sql.ErrNoRows {
		http.Error(w, "user if not defined", http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	if resp != nil {
		var grpcParam = &auth.ParamUser{
			Email: recoveryPasswordData.Email,
			Phone: recoveryPasswordData.Phone,
			Id:    0,
		}

		getParamForQuery := helpers.GetParamUser(grpcParam)
		ip := helpers.GetClientIP(r)

		// ключ содержит ip адрес и юзера(контактная информация), чтобы два пользователя с одного ip могли отпралять запрос независимо друг от друга
		key := fmt.Sprintf("send:code-email:%s:%s", getParamForQuery, ip)

		attempts, err := rdb.Incr(ctx, key).Result()
		if err != nil && err != redis.Nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if attempts == 1 {
			rdb.Expire(ctx, key, time.Minute)
		}

		if attempts > 1 {
			time.Sleep(500 * time.Millisecond) // защита от брутфорс атак
			http.Error(w, "Too many attempts. Try again later.", http.StatusTooManyRequests)
			return
		}

		if email != "" {
			code := helpers.SixRandomNumbers()

			queryUpdateCode := "UPDATE users SET code_recovery = $1, code_recovery_expires_at = NOW() + INTERVAL '10 minutes' WHERE id = $2"

			rows, err := db.ExecContext(ctx, queryUpdateCode, code, userId)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			rowsUpdated, _ := rows.RowsAffected()

			if rowsUpdated != 1 {
				http.Error(w, "User is not defined", http.StatusInternalServerError)
				return
			}

			from := resp.SenderEmail
			password := resp.Password

			to := []string{email}

			msg := []byte(
				"Subject: Восстановление пароля\r\n" +
					"Content-Type: text/plain; charset=\"UTF-8\"\r\n" +
					"\r\n" +
					fmt.Sprintf("Ваш код: %s", code),
			)

			auth := smtp.PlainAuth("", from, password, "smtp.gmail.com")

			err = smtp.SendMail(
				"smtp.gmail.com:587",
				auth,
				from,
				to,
				msg,
			)

			if err != nil {
				fmt.Println(err)
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "code sent successfully",
			})

		}
	} else {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
}
