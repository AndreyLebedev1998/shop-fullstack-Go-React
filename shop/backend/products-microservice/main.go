package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"products-microservice/cors"
	"products-microservice/endpoints/categories"
	"products-microservice/endpoints/products"
	grpc_package "products-microservice/gRPC"
	"products-microservice/models"
	productstats "products-microservice/product_stats"
	"time"

	"github.com/AndreyLebedev1998/auth-grpc"
	order "github.com/AndreyLebedev1998/shop-gRPC-orders"
	product "github.com/AndreyLebedev1998/shop-gRPC-product"
	"github.com/segmentio/kafka-go"
	httpSwagger "github.com/swaggo/http-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"

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
	topicName := "product_stats"
	groupId := "product_stats_group_id"
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

	defer connAuth.Close()

	clientAuth := auth.NewAuthServiceClient(connAuth)

	connOrders, err := grpc.NewClient("orders-microservice:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
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

	clientOrders := order.NewOrderServiceClient(connOrders)

	go func() {
		fmt.Println("gRPC server started on :50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{os.Getenv("KAFKA_BROKERS")},
		Topic:   topicName,
		GroupID: groupId,
	})

	go func() {
		for {
			msg, err := reader.ReadMessage(ctx)
			if err != nil {
				log.Println("Error:", err)
				continue
			}

			var stats []models.ProductStats

			if err := json.Unmarshal(msg.Value, &stats); err != nil {
				if errors.Is(err, context.Canceled) {
					log.Println("kafka reader stopped: context cancelled")
					return
				}
				log.Println("kafka read error:", err)
				continue
			}

			fmt.Println("msgKafka", stats)

			productstats.ProductStats(db, stats)
		}
	}()

	mux := http.NewServeMux()

	mux.Handle("/products-for-categories", cors.WithCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		products.GetAllProductsByCategoryId(w, r, db, rdb)
	})))

	mux.Handle("/products-for-subcategory", cors.WithCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		products.GetAllProductsBySubcategoryId(w, r, db, rdb)
	})))

	mux.Handle("/categories", cors.WithCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		categories.GetAllCategories(w, r, db, rdb)
	})))

	mux.Handle("/get-products-for-filters", cors.WithCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		products.GetProductsForFilters(w, r, db)
	})))

	mux.Handle("/get-initial-values-for-filter", cors.WithCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		products.GetInitialValuesForFilter(w, r, db)
	})))

	mux.Handle("/add-favorite-product-for-user", cors.WithCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		products.AddFavoriteProductForUser(w, r, db, clientAuth)
	})))

	mux.Handle("/get-favorite-products", cors.WithCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		products.GetFavoriteProducts(w, r, db, clientAuth)
	})))

	mux.Handle("/remove-favorite-products", cors.WithCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		products.RemoveFavoriteProduct(w, r, db, clientAuth)
	})))

	mux.Handle("/get-recommendations", cors.WithCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		products.GetRecommendationsForUser(w, r, db, rdb, clientAuth, clientOrders)
	})))

	mux.Handle("/find-product-by-symbols", cors.WithCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		products.FindProductBySymbols(w, r, db)
	})))

	mux.Handle("/sort-products", cors.WithCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		products.SortProductsForIndicator(w, r, db)
	})))

	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	http.ListenAndServe(":8080", mux)
}
