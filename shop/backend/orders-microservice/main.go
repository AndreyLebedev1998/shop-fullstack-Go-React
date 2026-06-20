package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"orders-microservice/cors"
	"orders-microservice/endpoints/orders"
	grpc_package "orders-microservice/gRPC"
	"os"
	"time"

	"github.com/AndreyLebedev1998/admin-grpc"
	auth "github.com/AndreyLebedev1998/auth-grpc"
	order "github.com/AndreyLebedev1998/shop-gRPC-orders"
	product "github.com/AndreyLebedev1998/shop-gRPC-product"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"

	_ "github.com/lib/pq"
)

func main() {
	connStr := os.Getenv("DB_CONN")
	var db *sql.DB

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}

	fmt.Println("✅ Подключено к PostgreSQL orders-microservice")

	var ctx = context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:        os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password:    os.Getenv("REDIS_PASSWORD"),
		DB:          0,
		MaxRetries:  3,
		DialTimeout: 3 * time.Second,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		panic(err)
	}

	fmt.Println("✅ Подключено к Redis orders-microservice")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("🔥 Работает! Сервер перезапустился!"))
	})

	mux := http.NewServeMux()

	conn, err := grpc.Dial("products-microservice:50051", grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff: backoff.Config{
				BaseDelay:  100 * time.Millisecond,
				MaxDelay:   5 * time.Second,
				Multiplier: 1.6,
				Jitter:     0.2,
			},

			MinConnectTimeout: 5 * time.Second,
		}),
	)

	if err != nil {
		log.Fatal("gRPC dial error:", err)
	}

	defer conn.Close()

	client := product.NewProductsServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()

	order.RegisterOrderServiceServer(grpcServer, &grpc_package.Server{
		DB:  db,
		RDB: rdb,
	})

	go func() {
		fmt.Println("gRPC server started on :50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()

	connAuth, err := grpc.Dial("authorization-microservice:50051", grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff: backoff.Config{
				BaseDelay:  100 * time.Millisecond,
				MaxDelay:   5 * time.Second,
				Multiplier: 1.6,
				Jitter:     0.2,
			},

			MinConnectTimeout: 5 * time.Second,
		}),
	)

	if err != nil {
		log.Fatal("gRPC dial error:", err)
	}

	defer conn.Close()

	clientAuth := auth.NewAuthServiceClient(connAuth)

	connAdmin, err := grpc.Dial("admin-microservice:50051", grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff: backoff.Config{
				BaseDelay:  100 * time.Millisecond,
				MaxDelay:   5 * time.Second,
				Multiplier: 1.6,
				Jitter:     0.2,
			},

			MinConnectTimeout: 5 * time.Second,
		}),
	)

	if err != nil {
		log.Fatal("gRPC dial error:", err)
	}

	defer conn.Close()

	clientAdmin := admin.NewAdminServiceClient(connAdmin)

	ctxAdmin, cancelAdmin := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancelAdmin()

	var bot *tgbotapi.BotAPI

	// Получение телегарм бота
	token, err := clientAdmin.GetTelegramToken(ctxAdmin, &admin.Empty{})
	if err != nil {
		fmt.Println(err)
	}

	if token != nil {
		bot, err = tgbotapi.NewBotAPI(token.Token)
		if err != nil {
			fmt.Println(err)
		}
	}

	mux.Handle("/create-order", cors.WithCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		orders.CreateOrder(w, r, db, client, clientAuth, bot)
	})))

	mux.Handle("/change-order", cors.WithCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		orders.ChangeOrder(w, r, db, client, clientAuth, bot)
	})))

	mux.Handle("/get-orders-by-user", cors.WithCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		orders.GetOrdersByParametr(w, r, db, rdb, client, clientAuth)
	})))

	mux.Handle("/get-orders-by-date-from-user", cors.WithCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		orders.GetOrdersOneDateByUser(w, r, db, rdb, client)
	})))

	mux.Handle("/get-orders-from-user-between-date", cors.WithCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		orders.GetOrdersByUserBetweenDate(w, r, db, rdb, client)
	})))

	mux.Handle("/get-orders-one-date", cors.WithCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		orders.GetOrdersOneDate(w, r, db, rdb, client)
	})))

	mux.Handle("/get-order-by-id", cors.WithCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		orders.GetOrderById(w, r, db, client)
	})))

	http.ListenAndServe(":8080", mux)
}
