package helpers

import (
	"products-microservice/constants"
	"products-microservice/models"
	"slices"

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

func getAllIndicators() []constants.Indicator {
	return []constants.Indicator{
		constants.IndicatorAlphabeticalOrder,
		constants.IndicatorCheaper,
		constants.IndicatorHigherRating,
		constants.IndicatorMoreExpensive,
	}
}

func IsValidIndicator(indicator constants.Indicator) bool {
	return slices.Contains(getAllIndicators(), indicator)
}

func GetSortForIndicator(indicator constants.Indicator) string {
	if indicator == "alphabet" {
		return "ORDER BY product_name "
	} else if indicator == "more_expensive" {
		return "ORDER BY price DESC "
	} else if indicator == "cheaper" {
		return "ORDER BY price ASC "
	} else {
		return ""
	}
}
