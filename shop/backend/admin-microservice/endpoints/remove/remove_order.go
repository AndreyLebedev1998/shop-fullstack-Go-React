package remove

import (
	"encoding/json"
	"net/http"
	"strconv"

	order "github.com/AndreyLebedev1998/shop-gRPC-orders"
	product "github.com/AndreyLebedev1998/shop-gRPC-product"
)

func RemoveOrder(w http.ResponseWriter, r *http.Request, clientProduct product.ProductsServiceClient, clientOrder order.OrderServiceClient) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	ctx := r.Context()

	orderIdStr := r.URL.Query().Get("id")
	orderId, err := strconv.Atoi(orderIdStr)
	if err != nil || orderId <= 0 {
		http.Error(w, "invalid order id", http.StatusBadRequest)
		return
	}

	respGet, err := clientOrder.GetOrderItems(ctx, &order.OrderId{
		Id: int64(orderId),
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var itemsForProduct []*product.UpdateProductQuantity

	for _, item := range respGet.Items {
		itemsForProduct = append(itemsForProduct, &product.UpdateProductQuantity{
			ProductId: item.ProductId,
			Quantity:  item.Quantity,
		})
	}

	respProd, err := clientProduct.UpdateProductQuantityByIds(ctx, &product.UpdateProductQuantityRequest{
		OldItems: itemsForProduct,
		NewItems: itemsForProduct,
		IsCreate: false,
		IsDelete: true,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if respProd != nil && respProd.Success {
		respRemove, err := clientOrder.RemoveOrder(ctx, &order.OrderId{
			Id: int64(orderId),
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(respRemove)
	}
}
