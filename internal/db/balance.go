package db

import (
	"context"
	"fmt"

	"github.com/lenarlenar/gomart/internal/models"
)

func (db *DataBase) CreateWithdrawal(orderID, userID string, amount float64) error {
	query := `
	INSERT INTO
		withdrawal_flow (order_id, user_id, amount)
	VALUES ($1, $2, $3)
	`
	_, err := db.DB.ExecContext(context.Background(), query, orderID, userID, amount)
	if err != nil {
		return fmt.Errorf("не удалось создать вывод средств: %w", err)
	}

	return nil
}

func (db *DataBase) FindWithdrawalFlow(userID string) (*[]models.WithdrawalFlowItemDB, error) {
	query := `
		SELECT
			order_id,
			amount,
			processed_at
		FROM
			withdrawal_flow
		WHERE
			user_id = $1
	`
	var result []models.WithdrawalFlowItemDB
	rows, err := db.DB.QueryContext(context.Background(), query, userID)
	if err != nil {
		return nil, fmt.Errorf("не удалось выполнить запрос вывода средств: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item models.WithdrawalFlowItemDB
		if err := rows.Scan(&item.OrderID, &item.Amount, &item.ProcessedAt); err != nil {
			return nil, fmt.Errorf("ошибка при сканировании строки вывода средств: %w", err)
		}
		result = append(result, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка после чтения строк вывода средств: %w", err)
	}

	return &result, nil
}

func (db *DataBase) FindAccrualFlow(userID string) (*[]models.AccrualFlowItemDB, error) {
	var result []models.AccrualFlowItemDB
	query := `
		SELECT
			af.order_id,
			af.amount,
			af.processed_at
		FROM
			accrual_flow af
		LEFT JOIN orders o ON af.order_id = o.id
		WHERE
			o.user_id = $1`
	rows, err := db.DB.QueryContext(context.Background(), query, userID)
	if err != nil {
		return nil, fmt.Errorf("не удалось выполнить запрос начислений: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var item models.AccrualFlowItemDB
		if err := rows.Scan(&item.OrderID, &item.Amount, &item.ProcessedAt); err != nil {
			return nil, fmt.Errorf("ошибка при сканировании строки начислений: %w", err)
		}
		result = append(result, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка после чтения строк начислений: %w", err)
	}

	return &result, nil
}
