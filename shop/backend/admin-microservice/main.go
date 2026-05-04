package main

import (
	"admin-microservice/cors"
	"admin-microservice/endpoints/change-products-and-categories"
	"admin-microservice/endpoints/create-products-and-categories"
	"admin-microservice/endpoints/email"
	"admin-microservice/endpoints/orders"
	"admin-microservice/endpoints/remove"
	"admin-microservice/endpoints/telegram"
	grpc_package "admin-microservice/gRPC"
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	admin "github.com/AndreyLebedev1998/admin-grpc"
	//auth "github.com/AndreyLebedev1998/auth-grpc"
	order "github.com/AndreyLebedev1998/shop-gRPC-orders"
	product "github.com/AndreyLebedev1998/shop-gRPC-product"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
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
	fmt.Println(db)

	var ctx = context.Background()

	fmt.Println("✅ Подключено к PostgreSQL admin-microservice")

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

	fmt.Println("✅ Подключено к Redis admin-microservice")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("🔥 Работает! Сервер перезапустился!"))
	})

	mux := http.NewServeMux()

	connProduct, err := grpc.Dial("products-microservice:50051", grpc.WithTransportCredentials(insecure.NewCredentials()),
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

	defer connProduct.Close()

	connOrders, err := grpc.Dial("orders-microservice:50051", grpc.WithTransportCredentials(insecure.NewCredentials()),
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

	defer connOrders.Close()

	/* connAuth, err := grpc.Dial("authorization-microservice:50051", grpc.WithTransportCredentials(insecure.NewCredentials()),
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

	defer connAuth.Close() */

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()

	admin.RegisterAdminServiceServer(grpcServer, &grpc_package.Server{
		DB:  db,
		RDB: rdb,
	})
	go func() {
		fmt.Println("gRPC server started on :50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()

	clientProduct := product.NewProductsServiceClient(connProduct)
	clientOrder := order.NewOrderServiceClient(connOrders)
	//clientAuth := auth.NewAuthServiceClient(connAuth)

	_, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	mux.Handle("/create-product", cors.WithCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		create.CreateProduct(w, r, clientProduct)
	})))

	mux.Handle("/create-category", cors.WithCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		create.CreateCategory(w, r, clientProduct)
	})))

	mux.Handle("/uploads-products-images/", http.StripPrefix(
		"/uploads-products-images/",
		http.FileServer(http.Dir("./uploads-products-images")),
	))

	mux.Handle("/change-status-order", cors.WithCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		orders.ChangeStatusOrder(w, r, clientOrder)
	})))

	mux.Handle("/change-product", cors.WithCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		change.ChangeProduct(w, r, clientProduct)
	})))

	mux.Handle("/change-category", cors.WithCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		change.ChangeCategory(w, r, clientProduct)
	})))

	mux.Handle("/change-status-order-paid", cors.WithCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		orders.ChangeStatusPaidOrder(w, r, clientOrder)
	})))

	mux.Handle("/remove-product", cors.WithCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		remove.RemoveProduct(w, r, clientProduct)
	})))

	mux.Handle("/remove-category", cors.WithCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		remove.RemoveCategory(w, r, clientProduct)
	})))

	mux.Handle("/remove-order", cors.WithCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		remove.RemoveOrder(w, r, clientProduct, clientOrder)
	})))

	mux.Handle("/add-token", cors.WithCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		telegram.AddTokenTg(w, r, db)
	})))

	mux.Handle("/remove-token", cors.WithCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		telegram.RemoveToken(w, r, db)
	})))

	mux.Handle("/change-email-conn-data", cors.WithCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		email.AddAccountDataForEmail(w, r, db)
	})))

	http.ListenAndServe(":8080", mux)
}
