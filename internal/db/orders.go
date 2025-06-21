package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/lenarlenar/gomart/internal/models"
	"github.com/lib/pq"
)

func (db *DataBase) CreateOrder(orderID string, userID string) error {
	_, err := db.DB.ExecContext(context.Background(), "INSERT INTO orders (id, user_id) VALUES ($1, $2)", orderID, userID)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return models.ErrDuplicateOrder
		}
		return fmt.Errorf("ошибка создания заказа: %w", err)
	}
	return nil
}

func (db *DataBase) FindOrder(orderID string) (*models.OrderDB, error) {
	order := &models.OrderDB{}
	query := "SELECT id, user_id, status, uploaded_at FROM orders WHERE id = $1"
	err := db.DB.QueryRowContext(context.Background(), query, orderID).
		Scan(&order.ID, &order.UserID, &order.Status, &order.UploadedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("ошибка поиска заказа: %w", err)
	}

	return order, nil
}

func (db *DataBase) FindOrdersWithAccrual(userID string) (*[]models.OrderWithAccrualDB, error) {

	query := `
		SELECT
			o.id,
			user_id,
			status,
			uploaded_at,
			SUM(coalesce(amount, 0))
		FROM
			orders o
			LEFT JOIN accrual_flow af ON o.id = af.order_id
		WHERE
			user_id = $1
		GROUP BY 
			o.id
	`

	var result []models.OrderWithAccrualDB

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	rows, err := db.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("ошибка поиска заказов с начислениями: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item models.OrderWithAccrualDB
		if err := rows.Scan(&item.ID, &item.UserID, &item.Status, &item.UploadedAt, &item.Accrual); err != nil {
			return nil, fmt.Errorf("ошибка обработки строки с заказом: %w", err)
		}
		result = append(result, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка итерации по строкам: %w", err)
	}

	return &result, nil
}
