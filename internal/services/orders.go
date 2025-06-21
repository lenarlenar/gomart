package services

import (
	"errors"
	"strconv"

	"github.com/lenarlenar/gomart/internal/interfaces"
	"github.com/lenarlenar/gomart/internal/models"
)

type OrdersService struct {
	OrdersStorage interfaces.OrdersStorage
}

func (s *OrdersService) Check(orderID string) bool {
	var sum int
	var alternate bool

	for i := len(orderID) - 1; i >= 0; i-- {
		digit, err := strconv.Atoi(string(orderID[i]))
		if err != nil {
			return false
		}
		if alternate {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
		alternate = !alternate
	}
	return sum%10 == 0
}

func (s *OrdersService) CreateOrder(orderID, userID string) error {
	if err := s.OrdersStorage.CreateOrder(orderID, userID); err != nil {
		if !errors.Is(err, models.ErrDuplicateOrder) {
			return err
		}

		order, errOrder := s.OrdersStorage.FindOrder(orderID)
		if errOrder != nil {
			return errOrder
		}

		if order.UserID == userID {
			return models.ErrDuplicateOrderByUser
		}

		return models.ErrDuplicateOrder
	}
	return nil
}

func (s *OrdersService) GetOrders(userID string) ([]models.Order, error) {

	orders, err := s.OrdersStorage.FindOrdersWithAccrual(userID)
	if err != nil {
		return []models.Order{}, err
	}

	if orders == nil {
		return []models.Order{}, nil
	}

	result := make([]models.Order, len(*orders))
	for i, order := range *orders {
		accrual := order.Accrual
		result[i] = models.Order{
			ID:         order.ID,
			Status:     order.Status.OrderStatus,
			UploadedAt: order.UploadedAt,
			Accrual:    &accrual,
		}
	}

	return result, nil
}
