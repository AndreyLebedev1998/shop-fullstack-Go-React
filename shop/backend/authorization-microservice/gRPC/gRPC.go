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

func (s *Server) AddTelegramToken(ctx context.Context, token *auth.NewTelegramToken) (*auth.TelegramToken, error) {
	query := "INSERT INTO telegram_tokens (id, token) VALUES (1, $1) ON CONFLICT (id) DO UPDATE SET token = EXCLUDED.token RETURNING id"

	var id int64

	err := s.DB.QueryRowContext(ctx, query, token.Token).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("server error: %v", err)
	}

	return &auth.TelegramToken{
		Token: token.Token,
		Id:    id,
	}, nil
}

func (s *Server) GetTelegramToken(ctx context.Context, empty *auth.Empty) (*auth.TelegramToken, error) {
	query := "SELECT id, token FROM telegram_tokens WHERE id = 1"

	var token auth.TelegramToken

	err := s.DB.QueryRowContext(ctx, query).Scan(&token.Id, &token.Token)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("Telegram token is not defined")
	}
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	return &auth.TelegramToken{
		Id:    token.Id,
		Token: token.Token,
	}, nil
}

func (s *Server) RemoveTelegramToken(ctx context.Context, empty *auth.Empty) (*auth.Message, error) {
	query := "DELETE FROM telegram_tokens WHERE id = 1"

	resp, err := s.DB.ExecContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	rowsDeleted, _ := resp.RowsAffected()

	if rowsDeleted != 1 {
		return nil, fmt.Errorf("Token is not defined")
	}

	return &auth.Message{
		Message: "Token deleted successfully",
	}, err
}
