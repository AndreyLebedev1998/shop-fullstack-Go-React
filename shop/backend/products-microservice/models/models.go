package models

type Product struct {
	Id                   int     `json:"id"`
	ProductName          string  `json:"product_name"`
	Price                float64 `json:"price"`
	CategoryId           int     `json:"category_id"`
	ImageUrl             *string `json:"image_url"`
	AvailabilityOfPieces int     `json:"availability_of_pieces"`
	SubcategoryId        int     `json:"subcategory_id"`
	CategoryName         string  `json:"category_name"`
}

type Category struct {
	Id           int    `json:"id"`
	CategoryName string `json:"category_name"`
}

type CategoryWithSubcategories struct {
	Id            int           `json:"id"`
	CategoryName  string        `json:"category_name"`
	Subcategories []Subcategory `json:"subcategories"`
}

type Subcategory struct {
	CategoryId    int    `json:"category_id"`
	CategoryName  string `json:"category_name"`
	SubcategoryId int    `json:"subcategory_id"`
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

type ProductsSubcategories struct {
	Id                   int     `json:"id"`
	ProductName          string  `json:"product_name"`
	Price                float64 `json:"price"`
	CategoryId           int     `json:"category_id"`
	ImageUrl             *string `json:"image_url"`
	AvailabilityOfPieces int     `json:"availability_of_pieces"`
	SubcategoryId        int     `json:"subcategory_id"`
	SubcategoryName      string  `json:"category_name"`
}

type InitialValuesForFilter struct {
	MinPrice float64 `json:"min_price"`
	MaxPrice float64 `json:"max_price"`
}

type ProductId struct {
	ProductId int `json:"product_id"`
}

type AllFavoriteProduct struct {
	Id                   int     `json:"id"`
	ProductName          string  `json:"product_name"`
	Price                float64 `json:"price"`
	CategoryId           int     `json:"category_id"`
	ImageUrl             *string `json:"image_url"`
	AvailabilityOfPieces int     `json:"availability_of_pieces"`
	SubcategoryId        int     `json:"subcategory_id"`
	UserId               int     `json:"user_id"`
}

type FavoriteProduct struct {
	Id        int    `json:"id"`
	ProductId int    `json:"product_id"`
	UserId    int    `json:"user_id"`
	CreatedAt string `json:"created_at"`
}

type ProductStats struct {
	ProductId     int `json:"product_id"`
	PurchaseCount int `json:"purchase_count"`
}

type RecommendationProduct struct {
	Id                   int     `json:"id"`
	ProductName          string  `json:"product_name"`
	Price                float64 `json:"price"`
	CategoryId           int     `json:"category_id"`
	ImageUrl             *string `json:"image_url"`
	AvailabilityOfPieces int     `json:"availability_of_pieces"`
	SubcategoryId        int     `json:"subcategory_id"`
	CategoryName         string  `json:"category_name"`
	Rating               string  `json:"rating"`
}
