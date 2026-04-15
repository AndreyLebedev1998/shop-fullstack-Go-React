package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"products-microservice/cors"
	"products-microservice/endpoints/categories"
	"products-microservice/endpoints/products"
	grpc_package "products-microservice/gRPC"
	"time"

	product "github.com/AndreyLebedev1998/shop-gRPC-product"
	httpSwagger "github.com/swaggo/http-swagger"
	"google.golang.org/grpc"

	"github.com/redis/go-redis/v9"

	_ "products-microservice/docs"

	_ "github.com/lib/pq"
)

// @title Vanilla Go API
// @version 1.0
// @description prodcuts-microservice
// @host localhost:8090
// @BasePath /

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

	fmt.Println("✅ Подключено к PostgreSQL products-microservice")

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

	fmt.Println("✅ Подключено к Redis products-microservice")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("🔥 Работает! Сервер перезапустился!"))
	})

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()

	product.RegisterProductsServiceServer(grpcServer, &grpc_package.Server{
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

	mux.Handle("/products", cors.WithCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		products.GetAllProductsByCategoryId(w, r, db, rdb)
	})))

	mux.Handle("/categories", cors.WithCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		categories.GetAllCategories(w, r, db, rdb)
	})))

	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	http.ListenAndServe(":8080", mux)
}
