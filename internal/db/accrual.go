package db

import (
	"context"
	"fmt"

	"github.com/lenarlenar/gomart/internal/models"
)

func (db *DataBase) Create(orderID string, amount float64) error {
	_, err := db.DB.ExecContext(context.Background(), "INSERT INTO accrual_flow (order_id, amount) VALUES ($1, $2)", orderID, amount)
	if err != nil {
		return fmt.Errorf("не удалось создать начисление: %w", err)
	}

	return nil
}

func (db *DataBase) UpdateStatus(orderID string, status models.OrderStatusDB) error {
	_, err := db.DB.ExecContext(context.Background(), "UPDATE orders SET status = $2 WHERE id = $1", orderID, status.OrderStatus)
	if err != nil {
		return fmt.Errorf("ошибка обновления статуса заказа: %w", err)
	}
	return nil
}

func (db *DataBase) FindAllUnprocessedOrders() (*[]models.OrderDB, error) {
	query := `
	SELECT
		id,
		user_id,
		status,
		uploaded_at
	FROM
		orders
	WHERE
		status NOT IN ('INVALID', 'PROCESSED')
`
	var result []models.OrderDB

	rows, err := db.DB.QueryContext(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("ошибка поиска необработанных заказов: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item models.OrderDB
		if err := rows.Scan(&item.ID, &item.UserID, &item.Status, &item.UploadedAt); err != nil {
			return nil, fmt.Errorf("ошибка обработки строки с заказом: %w", err)
		}
		result = append(result, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка итерации по строкам: %w", err)
	}

	return &result, nil
}
