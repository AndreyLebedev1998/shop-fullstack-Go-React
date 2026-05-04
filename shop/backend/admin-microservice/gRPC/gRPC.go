package grpc_package

import (
	"context"
	"database/sql"
	"fmt"

	admin "github.com/AndreyLebedev1998/admin-grpc"
	"github.com/redis/go-redis/v9"
)

type Server struct {
	admin.UnimplementedAdminServiceServer
	DB  *sql.DB
	RDB *redis.Client
}

func (s *Server) GetTelegramToken(ctx context.Context, empty *admin.Empty) (*admin.TelegramToken, error) {
	query := "SELECT id, token FROM telegram_tokens WHERE id = 1"

	var token admin.TelegramToken

	err := s.DB.QueryRowContext(ctx, query).Scan(&token.Id, &token.Token)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("Telegram token is not defined")
	}
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	return &admin.TelegramToken{
		Id:    token.Id,
		Token: token.Token,
	}, nil
}

func (s *Server) GetDataFromConnectionEmail(ctx context.Context, empty *admin.Empty) (*admin.DataFromConnEmail, error) {
	query := "SELECT sender_email, password FROM email_conn WHERE id = 1"

	var grpcData admin.DataFromConnEmail

	err := s.DB.QueryRowContext(ctx, query).Scan(&grpcData.SenderEmail, &grpcData.Password)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("%v", err)
	}
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	return &grpcData, nil
}
