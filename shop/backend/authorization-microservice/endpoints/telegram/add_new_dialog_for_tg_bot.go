package telegram

import (
	"authorization-microservice/helpers"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/AndreyLebedev1998/auth-grpc"
)

func NewDialogTgBot(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	var userIdStr = r.URL.Query().Get("id")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil && userIdStr != "" {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}
	phone := r.URL.Query().Get("phone")
	email := r.URL.Query().Get("email")

	var chatId *int

	queryGetChatId := "SELECT chat_id_telegram FROM users "

	if userId != 0 && userId > 0 {
		queryGetChatId += "WHERE id = $1"
	} else if phone != "" {
		queryGetChatId += "WHERE phone = $1"
	} else if email != "" {
		queryGetChatId += "WHERE email = $1"
	} else {
		http.Error(w, "phone, user_id, and email is empty", http.StatusBadRequest)
		return
	}

	var grpcParam = &auth.ParamUser{
		Email: email,
		Phone: phone,
		Id:    int64(userId),
	}

	getParamForQuery := helpers.GetParamUser(grpcParam)

	if getParamForQuery == nil {
		fmt.Println(getParamForQuery)
		http.Error(w, "phone, user_id, and email is empty", http.StatusBadRequest)
		return
	}

	err = db.QueryRowContext(ctx, queryGetChatId, getParamForQuery).Scan(&chatId)

	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	if chatId == nil {
		token := helpers.RandomString()

		queryUpdateToken := "UPDATE users SET temporary_token_tg = $1, telegram_token_expires_at = NOW() + INTERVAL '10 minutes' WHERE id = $2"

		rows, err := db.ExecContext(ctx, queryUpdateToken, token, userId)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Server error", http.StatusBadRequest)
			return
		}

		rowsUpdated, _ := rows.RowsAffected()

		if rowsUpdated != 1 {
			http.Error(w, "user is not defined", http.StatusBadRequest)
			return
		}

		queryGetUser := "SELECT id FROM users WHERE temporary_token_tg = $1 AND telegram_token_expires_at > NOW()"

		rows, err = db.ExecContext(ctx, queryGetUser, token)

		if err != nil {
			fmt.Println(err)
			http.Error(w, "Server error", http.StatusBadRequest)
			return
		}

		rowsGet, _ := rows.RowsAffected()

		if rowsGet != 1 {
			http.Error(w, "user is not defined or the token has expired", http.StatusBadRequest)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{
			"link": fmt.Sprintf("https://t.me/orders_shop_go_bot?start=%s", token),
		})
	} else {
		http.Error(w, "Похоже вы уже привязали бота к номеру телефона, почте или к аккаунту", http.StatusBadRequest)
		return
	}
}
