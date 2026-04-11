package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"orders-microservice/cors"
	"orders-microservice/endpoints/orders"
	"os"
	"time"

	"github.com/redis/go-redis/v9"

	_ "github.com/lib/pq"
)

func main() {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	var db *sql.DB

	psqlInfo := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s sslmode=disable",
		host, user, password, dbname,
	)

	db, err := sql.Open("postgres", psqlInfo)
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

	mux.Handle("/create-order", cors.WithCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		orders.CreateOrder(w, r, db)
	})))

	mux.Handle("/change-order", cors.WithCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		orders.ChangeOrder(w, r, db)
	})))

	mux.Handle("/get-orders-by-parametr", cors.WithCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		orders.GetOrdersByParametr(w, r, db)
	})))

	mux.Handle("/get-orders-by-date-from-user", cors.WithCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		orders.GetOrdersOneDateByUser(w, r, db)
	})))

	http.ListenAndServe(":8080", mux)
}
