package models

type Order struct {
	Id         int         `json:"id"`
	UserId     *int        `json:"user_id"`
	Email      *string     `json:"email"`
	Phone      *string     `json:"phone"`
	Status     *string     `json:"status"`
	TotalPrice float64     `json:"total_price"`
	CreatedAt  string      `json:"created_at"`
	UpdatedAt  string      `json:"updated_at"`
	OrderItems []OrderItem `json:"order_items"`
}

type NewOrder struct {
	UserId     *int        `json:"user_id"`
	Email      *string     `json:"email"`
	Phone      *string     `json:"phone"`
	Status     *string     `json:"status"`
	TotalPrice float64     `json:"total_price"`
	OrderItems []OrderItem `json:"order_items"`
}

type OrderItem struct {
	ProductId int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

type NewOrderId struct {
	Id int `json:"id"`
}

type FullOrder struct {
	OrderId    int        `json:"order_id"`
	UserId     *int       `json:"user_id"`
	Email      *string    `json:"email"`
	Phone      *string    `json:"phone"`
	Status     string     `json:"status"`
	TotalPrice string     `json:"total_price"`
	CreatedAt  string     `json:"created_at"`
	UpdatedAt  string     `json:"updated_at"`
	Products   []Products `json:"products"`
}

type Products struct {
	OrderItemId  string  `json:"order_item_id"`
	ProductId    int     `json:"product_id"`
	Quantity     int     `json:"quantity"`
	Price        float64 `json:"price"`
	ProductName  string  `json:"product_name"`
	CategoryId   int     `json:"category_id"`
	CategoryName string  `json:"category_name"`
	ImageUrl     *string `json:"image_url"`
}

type ProblemProducts struct {
	ProductId            int    `json:"product_id"`
	AvailabilityOfPieces int    `json:"availability_of_pieces"`
	ProductName          string `json:"product_name"`
}

type ProductsCheck struct {
	ProductId            int    `json:"product_id"`
	AvailabilityOfPieces int    `json:"availability_of_pieces"`
	ProductName          string `json:"product_name"`
}
