package grpc_package

import (
	"authorization-microservice/helpers"
	"context"
	"database/sql"
	"fmt"

	auth "github.com/AndreyLebedev1998/auth-grpc"
	"github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

type Server struct {
	auth.UnimplementedAuthServiceServer
	DB  *sql.DB
	RDB *redis.Client
}

func (s *Server) UpdateChatIdTgUser(ctx context.Context, chatAndToken *auth.ChatIdTgUser) (*auth.MessageForUpdateTgUser, error) {
	query := "UPDATE users SET chat_id_telegram = $1 WHERE temporary_token_tg = $2 AND telegram_token_expires_at > NOW() RETURNING id"
	var userId int

	err := s.DB.QueryRowContext(ctx, query, chatAndToken.ChatId, chatAndToken.Token).Scan(&userId)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("The token was not found or has expired")
	}

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return nil, fmt.Errorf("Похоже вы уже привязали бота к номеру телефона, почте или к аккаунту")
		}
		return nil, fmt.Errorf("Server error")
	}

	return &auth.MessageForUpdateTgUser{
		Msg: "Chat ID has been successfully linked to the user",
	}, nil
}

func (s *Server) GetChatIdForUser(ctx context.Context, paramUser *auth.ParamUser) (*auth.ChatId, error) {
	query := "SELECT chat_id_telegram FROM users "

	if paramUser.Id != 0 {
		query += "WHERE id = $1"
	} else if paramUser.Email != "" {
		query += "WHERE email = $1"
	} else if paramUser.Phone != "" {
		query += " WHERE phone = $1"
	}

	var chatId int64

	var userParamForQuery = helpers.GetParamUser(paramUser)

	if userParamForQuery != nil {
		err := s.DB.QueryRowContext(ctx, query, userParamForQuery).Scan(&chatId)
		if err == sql.ErrNoRows {
			fmt.Println(err)
			return nil, fmt.Errorf("user not found")
		}

		if err != nil {
			fmt.Println(err)
			return nil, fmt.Errorf("db error: %v", err)
		}

		return &auth.ChatId{
			ChatId: chatId,
		}, nil
	} else {
		return nil, fmt.Errorf("phone, user_id, and email is empty")
	}
}
