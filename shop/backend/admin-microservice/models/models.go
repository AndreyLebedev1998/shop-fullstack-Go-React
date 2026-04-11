package models

type NewProduct struct {
	ProductName string  `json:"product_name"`
	Price       float64 `json:"price"`
	CategoryId  int     `json:"category_id"`
	ImageUrl    *string `json:"image_url"`
}

type Product struct {
	Id          int     `json:"id"`
	ProductName string  `json:"product_name"`
	Price       float64 `json:"price"`
	CategoryId  int     `json:"category_id"`
	ImageUrl    *string `json:"image_url"`
}

type NewCategory struct {
	CategoryName string `json:"category_name"`
}

type Category struct {
	Id           int    `json:"id"`
	CategoryName string `json:"category_name"`
}

type OrderStatus struct {
	Status string `json:"status"`
}

type OrderStatusPaid struct {
	StatusPaid string `json:"status_paid"`
}
