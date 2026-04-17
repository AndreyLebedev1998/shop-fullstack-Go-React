package grpcFunctionsQuery

import (
	"database/sql"
	"orders-microservice/models"
)

func ChangeOrderStatus(db *sql.DB, orderStatus models.OrderStatus) (models.OrderStatusResponse, error) {
	query := "UPDATE orders SET status = $1 WHERE id = $2"

	rows, err := db.Exec(query, orderStatus.Status, orderStatus.Id)
	if err != nil {
		message := models.OrderStatusResponse{
			Response: false,
			Status:   "update status failed",
			Id:       orderStatus.Id,
		}
		return message, err
	}

	updatedRows, _ := rows.RowsAffected()

	if updatedRows == 0 {
		message := models.OrderStatusResponse{
			Response: false,
			Status:   "order is not defined",
			Id:       orderStatus.Id,
		}
		return message, nil
	}

	message := models.OrderStatusResponse{
		Response: true,
		Status:   "order status successfully updated",
		Id:       orderStatus.Id,
	}
	return message, nil
}
