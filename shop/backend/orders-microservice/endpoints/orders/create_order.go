package orders

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"orders-microservice/helpers"
	"orders-microservice/models"

	"github.com/AndreyLebedev1998/auth-grpc"
	product "github.com/AndreyLebedev1998/shop-gRPC-product"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func markOrderFailed(ctx context.Context, db *sql.DB, orderId int) {
	_, _ = db.ExecContext(ctx, "UPDATE orders SET status = $1 WHERE id = $2", "failed", orderId)
}

func CreateOrder(w http.ResponseWriter, r *http.Request, db *sql.DB, client product.ProductsServiceClient, clientAuth auth.AuthServiceClient, bot *tgbotapi.BotAPI) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var order models.NewOrder
	ctx := r.Context()

	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// смотрим, чтобы были проудкты
	if len(order.OrderItems) == 0 {
		http.Error(w, "Order must contain items", http.StatusBadRequest)
		return
	}

	var respUser *auth.UserInfo
	var errUserInfo error

	respUser, errUserInfo = clientAuth.GetUserFromContactInfo(ctx, &auth.ContactInfo{
		Email:  helpers.StringOrEmpty(order.Email),
		Phone:  helpers.StringOrEmpty(order.Phone),
		UserId: helpers.Int64OrZero(order.UserId),
	})

	fmt.Println(respUser)

	if errUserInfo != nil {
		fmt.Printf("%v\n", errUserInfo)
		return
	}

	if respUser != nil {
		userId := int(respUser.UserId)
		order.Email = &respUser.Email
		order.UserId = &userId
		order.Phone = &respUser.Phone
	}

	var product_ids []int

	// собираем id проудктов для их получения через сервис проудктов
	for _, p := range order.OrderItems {
		product_ids = append(product_ids, p.ProductId)
	}

	// запрос на поулчение продуктов
	resp, err := client.GetProductsByIds(ctx, &product.GetProductsRequest{
		ProductIds: helpers.ConvertIntToInt64(product_ids),
	})

	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	// проставление дефолтного статуса, если он отсутствует
	if order.Status == nil || *order.Status != "pending" {
		val := "pending"
		order.Status = &val
	}

	// запрос на создание заказа
	queryOrder := `
		INSERT INTO orders (user_id, email, phone, status, total_price)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	// запрос на создание элементов заказа
	queryOrderItems := `
		INSERT INTO order_items (order_id, product_id, quantity, price)
		VALUES ($1, $2, $3, $4)
	`

	var productsFromOrder []models.OrderItem      // продукты в заказе
	var problemProducts []models.ProblemProducts  // проблемные продукты(если пользователь заказал больше, чем есть ан складе)
	var productsInOrder []models.ProductsForOrder // тоже продукты в заказе
	for _, product := range order.OrderItems {
		productsFromOrder = append(productsFromOrder, product)
	}

	for i, product := range resp.Products {
		if product.AvailabilityOfPieces < int64(productsFromOrder[i].Quantity) {
			// сравниваем колличество продуктов на складе с тем, что есть в заказе
			// если в заказе больше, то добавляем продукт в слайс проблемных продуктов
			var problemProduct models.ProblemProducts
			problemProduct.AvailabilityOfPieces = int(product.AvailabilityOfPieces)
			problemProduct.ImageUrl = &product.ImageUrl
			problemProduct.ProductId = int(product.Id)
			problemProduct.ProductName = product.ProductName
			problemProducts = append(problemProducts, models.ProblemProducts(problemProduct))
		}

		var productInOrder models.ProductsForOrder

		productInOrder.Id = int(product.Id)
		productInOrder.ProductName = product.ProductName
		productInOrder.CategoryId = int(product.CategoryId)
		productInOrder.CategoryName = product.CategoryName
		productInOrder.ImageUrl = &product.ImageUrl
		productInOrder.Quantity = int(product.AvailabilityOfPieces)
		productInOrder.Price = product.Price

		// добавляем продкт в слайс (продукты в заказе)
		productsInOrder = append(productsInOrder, productInOrder)
	}

	// если есть хоть один проблемный продукт, то заказ не создаем, а возврааем эти проудкты с указанием остатка
	if len(problemProducts) > 0 {
		problemProductsMsg := map[string]interface{}{
			"message":          "Sorry, there were no products in stock",
			"problem_products": problemProducts,
		}
		json.NewEncoder(w).Encode(problemProductsMsg)
		return
	}

	// считаем конечную цену заказа
	totalPrice := 0.0
	for _, p := range productsInOrder {
		for _, product := range order.OrderItems {
			if int(p.Id) == int(product.ProductId) {
				totalPrice += (p.Price * float64(product.Quantity))
				continue
			}
		}
	}

	// объявление транзакции
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer func() {
		_ = tx.Rollback()
	}()

	var newOrderId models.NewOrderId

	// создаем заказ
	err = tx.QueryRowContext(ctx, queryOrder,
		order.UserId,
		order.Email,
		order.Phone,
		order.Status,
		0,
	).Scan(&newOrderId.Id)

	if err != nil {
		http.Error(w, "Error creating order", http.StatusInternalServerError)
		return
	}

	// обновляем конечную цену
	_, err = tx.ExecContext(ctx, `
		UPDATE orders SET total_price = $1 WHERE id = $2
	`, totalPrice, newOrderId.Id)

	if err != nil {
		http.Error(w, "Failed to update total price", http.StatusInternalServerError)
		return
	}

	// добавляем в базе элементы заказа
	for _, item := range order.OrderItems {
		for _, p := range productsInOrder {
			// проходимся по продуктам в заказе  и по слайсу проудктов в заказе, на каждой итерации сравниваем id и добавляем элемент заказа в базу
			if item.ProductId == p.Id {
				_, err = tx.ExecContext(ctx, queryOrderItems,
					newOrderId.Id,
					item.ProductId,
					item.Quantity,
					p.Price,
				)

				if err != nil {
					http.Error(w, "Error creating order_items", http.StatusInternalServerError)
					return
				}
			}
		}
	}

	// формирование заказа, чтоыб вернуть его пользователю
	var query string = `SELECT orders.id as order_id, user_id, email, phone, status, total_price, created_at,
							product_id, quantity, price
							FROM orders
							JOIN order_items ON orders.id = order_items.order_id
							WHERE orders.id = $1`

	var products []models.Products
	var fullOrder models.FullOrder
	rows, err := tx.QueryContext(ctx, query, newOrderId.Id)
	if err != nil {
		http.Error(w, "order receiving error", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	for rows.Next() {
		var item models.Products
		err := rows.Scan(
			&fullOrder.OrderId,
			&fullOrder.UserId,
			&fullOrder.Email,
			&fullOrder.Phone,
			&fullOrder.Status,
			&fullOrder.TotalPrice,
			&fullOrder.CreatedAt,
			&item.ProductId,
			&item.Quantity,
			&item.Price,
		)
		if err != nil {
			http.Error(w, "Error reading row", http.StatusInternalServerError)
			return
		}
		products = append(products, item)
	}

	fullOrder.Products = products

	if err := tx.Commit(); err != nil {
		http.Error(w, "Commit failed", http.StatusInternalServerError)
		return
	}

	// создание слайса, для того, чтобы вычесть со склада количество продуктов заказанные пользователем
	grpcItems := make([]*product.UpdateProductQuantity, 0, len(productsFromOrder))

	// формирование слайса
	for _, p := range productsFromOrder {
		grpcItems = append(grpcItems, &product.UpdateProductQuantity{
			ProductId: int64(p.ProductId),
			Quantity:  int64(p.Quantity),
		})
	}

	// запрос на уменьшение кол-ва проудктов на складе заказанных пользователем
	response, err := client.UpdateProductQuantityByIds(ctx, &product.UpdateProductQuantityRequest{
		NewItems: grpcItems,
		OldItems: grpcItems,
		IsCreate: true,
		IsDelete: false,
	})

	// если ошибка, и мы не вычли со скалда проудкты заказанные пользователм, помечаем заказа как неудавшийся
	if err != nil {
		markOrderFailed(ctx, db, newOrderId.Id)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !response.Success {
		markOrderFailed(ctx, db, newOrderId.Id)
		http.Error(w, response.Message, http.StatusInternalServerError)
		return
	}

	var userEmail string
	if fullOrder.Email != nil {
		userEmail = *fullOrder.Email
	}

	var userId int64
	if fullOrder.UserId != nil {
		userId = int64(*fullOrder.UserId)
	}

	var userPhone string
	if fullOrder.Phone != nil {
		userPhone = *fullOrder.Phone
	}

	// получение чата id пользователя для отправки заказа в тг
	respAuth, err := clientAuth.GetChatIdForUser(ctx, &auth.ParamUser{
		Email: userEmail,
		Id:    userId,
		Phone: userPhone,
	})

	// отправка заказа пользователю в тг
	if err != nil {
		fmt.Println(err)
		fmt.Println("Ошибка отправки заказа в Telegram")
	} else {
		if bot != nil {
			orderTg := helpers.FormatOrderMessage(fullOrder, true)

			msg := tgbotapi.NewMessage(respAuth.ChatId, orderTg)
			msg.ParseMode = "HTML"
			_, err := bot.Send(msg)
			if err != nil {
				log.Println(err)
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	// возвращаем сформированный заказ пользователю
	json.NewEncoder(w).Encode(fullOrder)
}
