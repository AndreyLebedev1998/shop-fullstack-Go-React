package grpc_package

import (
	"authorization-microservice/constants"
	"authorization-microservice/helpers"
	"authorization-microservice/models"
	"context"
	"database/sql"
	"fmt"

	auth "github.com/AndreyLebedev1998/auth-grpc"
	"github.com/golang-jwt/jwt/v5"
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

func (s *Server) GetUserFromToken(ctx context.Context, token *auth.Token) (*auth.UserId, error) {
	claims := &models.Claims{}
	jwtToken, err := jwt.ParseWithClaims(token.Token, claims, func(t *jwt.Token) (interface{}, error) {
		return constants.JwtKey, nil
	})

	if err != nil || !jwtToken.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return &auth.UserId{
		UserId: int64(claims.UserID),
	}, nil
}

func (s *Server) GetUserFromContactInfo(ctx context.Context, req *auth.ContactInfo) (*auth.UserInfo, error) {
	query := `SELECT id, phone, email FROM users 
              WHERE ($1 != '' AND email = $1) 
              OR ($2 != '' AND phone = $2) 
              OR ($3 > 0 AND id = $3)
              LIMIT 1`

	var userInfo auth.UserInfo
	err := s.DB.QueryRowContext(ctx, query, req.Email, req.Phone, req.UserId).
		Scan(&userInfo.UserId, &userInfo.Phone, &userInfo.Email)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return &userInfo, nil
}
