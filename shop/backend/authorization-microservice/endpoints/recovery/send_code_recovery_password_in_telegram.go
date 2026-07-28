package recovery

import (
	"authorization-microservice/helpers"
	"authorization-microservice/interfaces"
	"authorization-microservice/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/AndreyLebedev1998/admin-grpc"
	"github.com/AndreyLebedev1998/auth-grpc"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/redis/go-redis/v9"
)

func SendCodeRecoveryPasswordInTelegram(w http.ResponseWriter, r *http.Request, db *sql.DB, rdb *redis.Client, clientAdmin admin.AdminServiceClient) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var recoveryPasswordData models.RecoveryPasswordData
	ctx := r.Context()
	var bot *tgbotapi.BotAPI

	if err := json.NewDecoder(r.Body).Decode(&recoveryPasswordData); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	var grpcParam = &auth.ParamUser{
		Email: recoveryPasswordData.Email,
		Phone: recoveryPasswordData.Phone,
		Id:    0,
	}

	token, err := clientAdmin.GetTelegramToken(ctx, &admin.Empty{})

	if err != nil || token == nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	getParamForQuery := helpers.GetParamUser(grpcParam)
	ip := helpers.GetClientIP(r)

	// ключ содержит ip адрес и юзера(контактная информация), чтобы два пользователя с одного ip могли отпралять запрос независимо друг от друга
	key := fmt.Sprintf("send:code-tg:%s:%s", getParamForQuery, ip)

	attempts, err := rdb.Incr(ctx, key).Result()
	if err != nil && err != redis.Nil {
		fmt.Println(err)
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

	if recoveryPasswordData.Email != "" || recoveryPasswordData.Phone != "" {
		var chatId *int64
		var userId int
		queryGetChatId := "SELECT id, chat_id_telegram FROM users "
		if recoveryPasswordData.Email != "" {
			queryGetChatId += "WHERE email = $1"
		} else {
			queryGetChatId += "WHERE phone = $1"
		}

		getParamForQuery := helpers.GetParamUser(grpcParam)

		err := db.QueryRowContext(ctx, queryGetChatId, getParamForQuery).Scan(&userId, &chatId)

		if err == sql.ErrNoRows {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "The Telegram bot is not linked to this user",
			})
			return
		}

		if chatId == nil || *chatId == 0 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "The Telegram bot is not linked to this user",
			})
			return
		}

		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if *chatId != 0 {
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

			bot, err = tgbotapi.NewBotAPI(token.Token)
			if err != nil {
				fmt.Println(err)
			}

			message := []byte("")
			var notifier interfaces.Notifier = interfaces.TgBot{
				ChatId: *chatId,
				Bot:    bot,
			}

			err = notifier.Send(message)

			if err != nil {
				log.Println(err)
				http.Error(w, "Failed to send message in Telegram", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "code sent successfully",
			})
		}
	} else {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
}
