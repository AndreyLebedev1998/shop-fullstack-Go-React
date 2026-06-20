package models

import "github.com/golang-jwt/jwt/v4"

type NewProduct struct {
	ProductName          string  `json:"product_name"`
	Price                float64 `json:"price"`
	CategoryId           int     `json:"category_id"`
	ImageUrl             *string `json:"image_url"`
	AvailabilityOfPieces int     `json:"availability_of_pieces"`
	SubcategoryId        int     `json:"subcategory_id"`
}

type Product struct {
	Id                   int     `json:"id"`
	ProductName          string  `json:"product_name"`
	Price                float64 `json:"price"`
	CategoryId           int     `json:"category_id"`
	ImageUrl             *string `json:"image_url"`
	AvailabilityOfPieces int     `json:"availability_of_pieces"`
	SubcategoryId        int     `json:"subcategory_id"`
}

type NewCategory struct {
	CategoryName string `json:"category_name"`
}

type Category struct {
	Id           int    `json:"id"`
	CategoryName string `json:"category_name"`
}

type NewSubcategory struct {
	CategoryName string `json:"category_name"`
	CategoryId   int    `json:"category_id"`
}

type Subcategory struct {
	Id           int    `json:"id"`
	CategoryName string `json:"category_name"`
	CategoryId   int    `json:"category_id"`
}

type OrderStatus struct {
	Status string
}

type OrderStatusPaid struct {
	StatusPaid string `json:"status_paid"`
}

type FullOrder struct {
	OrderId  int64      `json:"order_id"`
	Products []Products `json:"products"`
}

type Products struct {
	ProductId int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

type NewTelegramToken struct {
	Token string `json:"token"`
}

type EmailConnData struct {
	SenderEmail string `json:"sender_email"`
	Password    string `json:"password"`
}

type AuthAdmin struct {
	Nick     string `json:"nick"`
	Password string `json:"password"`
}

type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

type RegisterNewAdmin struct {
	Nick     string `json:"nick"`
	Password string `json:"password"`
}
