package models

import (
	"context"
	"fmt"
	"time"
)

type UserRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type User struct {
	ID       string
	Login    string
	HashPass string
}

type OrderStatus string

const (
	StatusNew        OrderStatus = "NEW"
	StatusProcessing OrderStatus = "PROCESSING"
	StatusInvalid    OrderStatus = "INVALID"
	StatusProcessed  OrderStatus = "PROCESSED"
)

type Order struct {
	ID         string      `json:"number"`
	Status     OrderStatus `json:"status"`
	Accrual    *float64    `json:"accrual,omitempty"`
	UploadedAt time.Time   `json:"uploaded_at"`
}

type OrderStatusDB struct {
	OrderStatus
}

func (o *OrderStatusDB) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("статус заказа должен быть строкой, а не %T", value)
	}

	o.OrderStatus = OrderStatus(string(bytes))
	return nil
}

type OrderDB struct {
	ID         string
	UserID     string
	Status     OrderStatusDB
	UploadedAt time.Time
}

type OrderWithAccrualDB struct {
	OrderDB
	Accrual float64
}

type Job func(ctx context.Context)

type Balance struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

type Withdrawal struct {
	ID  *string  `json:"order"`
	Sum *float64 `json:"sum"`
}

type WithdrawalFlowItem struct {
	OrderID     string    `json:"order"`
	Sum         float64   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}

type AccrualFlowItemDB struct {
	OrderID     string
	Amount      float64
	ProcessedAt time.Time
}

type WithdrawalFlowItemDB struct {
	OrderID     string
	Amount      float64
	ProcessedAt time.Time
}
