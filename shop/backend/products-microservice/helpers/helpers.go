package helpers

import (
	"products-microservice/models"

	product "github.com/AndreyLebedev1998/shop-gRPC-product"
)

func ConvertInt64ToInt(in []int64) []int {
	out := make([]int, len(in))
	for i, v := range in {
		out[i] = int(v)
	}
	return out
}

func ConvertIntToInt64(in []int) []int64 {
	out := make([]int64, len(in))
	for i, v := range in {
		out[i] = int64(v)
	}
	return out
}

func ConvertProductsToProto(products []models.ProductWithCategory) []*product.Product {
	result := make([]*product.Product, 0, len(products))

	for _, p := range products {
		result = append(result, &product.Product{
			Id:           int64(p.Id),
			ProductName:  p.ProductName,
			CategoryId:   int64(p.CategoryId),
			CategoryName: p.CategoryName,
			ImageUrl: func() string {
				if p.ImageUrl != nil {
					return *p.ImageUrl
				}
				return ""
			}(),
			AvailabilityOfPieces: int64(p.AvailabilityOfPieces),
			Price:                float64(p.Price),
		})
	}

	return result
}
