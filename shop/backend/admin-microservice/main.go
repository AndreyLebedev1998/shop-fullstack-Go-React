package main

import (
	"admin-microservice/cors"
	"admin-microservice/endpoints/create"
	"admin-microservice/endpoints/orders"
	"context"
	"log"
	"net/http"
	"time"

	order "github.com/AndreyLebedev1998/shop-gRPC-orders"
	product "github.com/AndreyLebedev1998/shop-gRPC-product"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

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

	clientProduct := product.NewProductsServiceClient(connProduct)
	clientOrder := order.NewOrderServiceClient(connOrders)

	_, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	mux.Handle("/create-product", cors.WithCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		create.CreateProduct(w, r, clientProduct)
	})))

	/* mux.Handle("/create-category", cors.WithCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		create.CreateCategory(w, r, db)
	}))) */

	mux.Handle("/uploads-products-images/", http.StripPrefix(
		"/uploads-products-images/",
		http.FileServer(http.Dir("./uploads-products-images")),
	))

	mux.Handle("/change-status-order", cors.WithCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		orders.ChangeStatusOrder(w, r, clientOrder)
	})))

	/* mux.Handle("/change-status-order-paid", cors.WithCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		orders.ChangeStatusPaidOrder(w, r, db)
	}))) */

	http.ListenAndServe(":8080", mux)
}
