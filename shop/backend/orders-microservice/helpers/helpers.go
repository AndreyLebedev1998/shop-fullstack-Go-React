package helpers

import (
	"database/sql"
	"fmt"
	"orders-microservice/models"
	"time"
)

func SqlQueryWithParam(params []models.ParamsForQuery) (string, []any) {
	query := ""

	var strParams []any
	for i, param := range params {
		if param.Value != "" {
			query += fmt.Sprintf(" AND %s = $%d", param.Column, i+1)
			strParams = append(strParams, param.Value)
		}
	}

	return query, strParams
}

func SqlQueryWithParamAndDate(params []models.ParamsForQuery) (string, []any) {
	query := `WHERE created_at >= $1 AND created_at < $2`

	var strParams []any
	for i, param := range params {
		if param.Value != "" {
			query += fmt.Sprintf(" AND %s = $%d", param.Column, i+3)
			strParams = append(strParams, param.Value)
		}
	}

	return query, strParams
}
func SqlQueryWithParamAndOneDate(params []models.ParamsForQuery) (string, []any) {
	query := "WHERE created_at >= $1 AND created_at < $2"

	var strParams []any
	for i, param := range params {
		if param.Value != "" {
			query += fmt.Sprintf(" AND %s = $%d", param.Column, i+3)
			strParams = append(strParams, param.Value)
		}
	}

	return query, strParams
}

func ForRowsAfterQuery(rows *sql.Rows, ordersMap map[int]*models.FullOrder, productIDsSet map[int]struct{}) error {
	for rows.Next() {
		var fullOrder models.FullOrder
		var product models.Products
		if err := rows.Scan(&fullOrder.OrderId, &fullOrder.UserId, &fullOrder.Email, &fullOrder.Phone,
			&fullOrder.Status, &fullOrder.TotalPrice, &fullOrder.CreatedAt, &fullOrder.UpdatedAt,
			&product.ProductId, &product.Quantity, &product.Price); err != nil {
			return fmt.Errorf("error scanning row: %w", err)
		}

		productIDsSet[product.ProductId] = struct{}{}

		if o, ok := ordersMap[fullOrder.OrderId]; !ok {
			fullOrder.Products = []models.Products{product}
			ordersMap[fullOrder.OrderId] = &fullOrder
		} else {
			o.Products = append(o.Products, product)
		}
	}
	return rows.Err()
}

func ConvertIntToInt64(in []int) []int64 {
	out := make([]int64, len(in))
	for i, v := range in {
		out[i] = int64(v)
	}
	return out
}

func derefString(s *string) string {
	if s == nil {
		return "-"
	}
	return *s
}

func derefInt(i *int) int {
	if i == nil {
		return 0
	}
	return *i
}

func formatTime(t string) string {
	parsed, err := time.Parse(time.RFC3339Nano, t)
	if err != nil {
		return t
	}

	return parsed.Format("02.01.2006 15:04:05")
}

func FormatOrderMessage(o models.FullOrder, isCreate bool) string {
	var productsText string
	var orderStatus string
	if isCreate {
		orderStatus = "Новый заказ"
	} else {
		orderStatus = "Заказ изменен"
	}

	for _, p := range o.Products {
		productsText += fmt.Sprintf(
			"• %s x%d — %.2f₽\n",
			p.ProductName,
			p.Quantity,
			p.Price,
		)
	}

	return fmt.Sprintf(
		"<b>🛒 %s #%d</b>\n\n"+
			"<b>Пользователь:</b> %v\n"+
			"<b>Email:</b> %v\n"+
			"<b>Phone:</b> %v\n\n"+
			"<b>Статус:</b> %s\n"+
			"<b>Сумма:</b> %.2f₽\n\n"+
			"<b>Товары:</b>\n%s\n"+
			"<i>Создан: %s</i>",
		orderStatus,
		o.OrderId,
		derefInt(o.UserId),
		derefString(o.Email),
		derefString(o.Phone),
		o.Status,
		o.TotalPrice,
		productsText,
		formatTime(o.CreatedAt),
	)
}

func StringOrEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func Int64OrZero(i *int) int64 {
	if i == nil {
		return 0
	}
	return int64(*i)
}
