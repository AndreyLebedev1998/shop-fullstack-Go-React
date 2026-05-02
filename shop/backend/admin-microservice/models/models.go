package models

type NewProduct struct {
	ProductName          string  `json:"product_name"`
	Price                float64 `json:"price"`
	CategoryId           int     `json:"category_id"`
	ImageUrl             *string `json:"image_url"`
	AvailabilityOfPieces int     `json:"availability_of_pieces"`
}

type Product struct {
	Id                   int     `json:"id"`
	ProductName          string  `json:"product_name"`
	Price                float64 `json:"price"`
	CategoryId           int     `json:"category_id"`
	ImageUrl             *string `json:"image_url"`
	AvailabilityOfPieces int     `json:"availability_of_pieces"`
}

type NewCategory struct {
	CategoryName string `json:"category_name"`
}

type Category struct {
	Id           int    `json:"id"`
	CategoryName string `json:"category_name"`
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
