package models

type Product struct {
	Id          int     `json:"id"`
	ProductName string  `json:"product_name"`
	Price       float64 `json:"price"`
	CategoryId  int     `json:"category_id"`
	ImageUrl    *string `json:"image_url"`
}

type Category struct {
	Id           int    `json:"id"`
	CategoryName string `json:"category_name"`
}

type ProductWithCategory struct {
	Id                   int
	ProductName          string
	Price                float64
	CategoryId           int
	CategoryName         string
	ImageUrl             *string
	AvailabilityOfPieces int
}

type UpdateProductForGRPC struct {
	ProductId int64 `json:"product_id"`
	Quantity  int64 `json:"quantity"`
}

type MessageUpdatedQuantityProducts struct {
	Success bool
	Message string
}
