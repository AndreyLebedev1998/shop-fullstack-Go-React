package main

import (
	"net/http"

	_ "github.com/lib/pq"
)

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("🔥 Работает! Сервер перезапустился!"))
	})

	mux := http.NewServeMux()

	/* mux.Handle("/create-product", cors.WithCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		create.CreateProduct(w, r, db)
	})))

	mux.Handle("/create-category", cors.WithCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		create.CreateCategory(w, r, db)
	})))

	mux.Handle("/uploads-products-images/", http.StripPrefix(
		"/uploads-products-images/",
		http.FileServer(http.Dir("./uploads-products-images")),
	))

	mux.Handle("/change-status-order", cors.WithCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		orders.ChangeStatusOrder(w, r, db)
	})))

	mux.Handle("/change-status-order-paid", cors.WithCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		orders.ChangeStatusPaidOrder(w, r, db)
	}))) */

	http.ListenAndServe(":8080", mux)
}
