package main

import (
	"authorization-microservice/cors"
	autentification "authorization-microservice/endpoints/autentification"
	"authorization-microservice/endpoints/telegram"
	grpc_package "authorization-microservice/gRPC"
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"

	auth "github.com/AndreyLebedev1998/auth-grpc"

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

	fmt.Println("✅ Подключено к PostgreSQL authorization-microservice")

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

	fmt.Println("✅ Подключено к Redis authorization-microservice")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("🔥 Работает! Сервер перезапустился!"))
	})

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()

	auth.RegisterAuthServiceServer(grpcServer, &grpc_package.Server{
		DB:  db,
		RDB: rdb,
	})

	go func() {
		fmt.Println("gRPC server started on :50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()

	mux := http.NewServeMux()

	mux.Handle("/registration", cors.WithCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		autentification.Register(w, r, db)
	})))

	mux.Handle("/authorization", cors.WithCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		autentification.Login(w, r, db)
	})))

	mux.Handle("/get-me", cors.WithCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		autentification.Me(w, r, db)
	})))

	mux.Handle("/add-new-tg-bot-dialog", cors.WithCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		telegram.NewDialogTgBot(w, r, db)
	})))

	http.ListenAndServe(":8080", mux)
}
