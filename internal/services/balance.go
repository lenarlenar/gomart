package services

import (
	"sort"

	"github.com/lenarlenar/gomart/internal/interfaces"
	"github.com/lenarlenar/gomart/internal/models"
)

type BalanceService struct {
	Storage interfaces.BalanceStorage
}

func (b *BalanceService) GetUserBalance(userID string) (models.Balance, error) {
	accrualFlow, err := b.Storage.FindAccrualFlow(userID)

	if err != nil {
		return models.Balance{}, err
	}

	withdrawalFlow, err := b.Storage.FindWithdrawalFlow(userID)

	if err != nil {
		return models.Balance{}, err
	}

	var current float64 = 0
	var withdrawn float64 = 0

	if accrualFlow != nil {
		for _, item := range *accrualFlow {
			current += item.Amount
		}
	}

	if withdrawalFlow != nil {
		for _, item := range *withdrawalFlow {
			withdrawn += item.Amount
		}
	}

	return models.Balance{Current: current - withdrawn, Withdrawn: withdrawn}, nil
}

func (b *BalanceService) CreateWithdrawal(orderID, userID string, amount float64) error {
	if err := b.Storage.CreateWithdrawal(orderID, userID, amount); err != nil {
		return err
	}

	return nil
}

func (b *BalanceService) GetWithdrawalFlow(userID string) ([]models.WithdrawalFlowItem, error) {
	withdrawalFlow, err := b.Storage.FindWithdrawalFlow(userID)

	if err != nil {
		return []models.WithdrawalFlowItem{}, err
	}

	if withdrawalFlow == nil {
		return []models.WithdrawalFlowItem{}, nil
	}

	result := make([]models.WithdrawalFlowItem, len(*withdrawalFlow))

	for i, item := range *withdrawalFlow {
		result[i] = models.WithdrawalFlowItem{
			OrderID:     item.OrderID,
			Sum:         item.Amount,
			ProcessedAt: item.ProcessedAt,
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].ProcessedAt.Before(result[j].ProcessedAt)
	})

	return result, nil
}
