package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/lenarlenar/gomart/internal/interfaces"
	"github.com/lenarlenar/gomart/internal/logger"
	"github.com/lenarlenar/gomart/internal/models"
	"go.uber.org/zap"
)

var (
	errNoOrder = errors.New("заказ не зарегистрирован")
	errServer  = errors.New("внутренняя ошибка сервера")
)

type accrualOrderStatus string

const (
	AccrualStatusRegistered accrualOrderStatus = "REGISTERED"
	AccrualStatusInvalid    accrualOrderStatus = "INVALID"
	AccrualStatusProcessing accrualOrderStatus = "PROCESSING"
	AccrualStatusProcessed  accrualOrderStatus = "PROCESSED"
)

type AccrualService struct {
	Storage          interfaces.AccrualStorage
	JobQueueService  interfaces.AccrualJobQueueService
	ExternalEndpoint string
}

type accrualDataResponse struct {
	ID      string             `json:"order"`
	Status  accrualOrderStatus `json:"status"`
	Accrual *float64           `json:"accrual,omitempty"`
}

const defaultRetryAfterDuration = 60

func (as *AccrualService) Calculate(orderID string) {
	as.JobQueueService.Enqueue(func(ctx context.Context) {
		data, retryAfter, err := fetchAccrualData(as.ExternalEndpoint, orderID)

		if err != nil {
			if errors.Is(err, errNoOrder) {
				logger.Log.Info("Заказ не зарегистрирован", zap.String("orderID", orderID))
				return
			}

			logger.Log.Error("Не удалось получить данные начислений", zap.Error(err))
			return
		}

		if retryAfter > 0 {
			logger.Log.Info("Получен Retry-After", zap.Int("retryAfter", retryAfter), zap.String("orderID", orderID))
			as.JobQueueService.PauseAndResume(time.Second * time.Duration(retryAfter))
			as.JobQueueService.Enqueue(func(ctx context.Context) {
				as.Calculate(orderID)
			})
			logger.Log.Info("Добавлено новое задание после паузы", zap.Int("retryAfter", retryAfter), zap.String("orderID", orderID))
			return
		}

		logger.Log.Info("Получены данные начислений",
			zap.String("orderID", orderID),
			zap.String("status", string(data.Status)),
		)

		switch data.Status {
		case AccrualStatusRegistered:
			as.JobQueueService.ScheduleJob(func(ctx context.Context) {
				as.Calculate(orderID)
			}, time.Minute)
			logger.Log.Info("Добавлено запланированное задание", zap.String("orderID", orderID))

		case AccrualStatusProcessed, AccrualStatusProcessing, AccrualStatusInvalid:
			err := as.Storage.UpdateStatus(orderID, models.OrderStatusDB{OrderStatus: models.OrderStatus(data.Status)})
			if err != nil {
				logger.Log.Error("Не удалось обновить статус заказа", zap.Error(err))
				return
			}

			logger.Log.Info("Статус заказа обновлен",
				zap.String("orderID", orderID),
				zap.String("status", string(data.Status)),
			)

			if data.Accrual != nil {
				err := as.Storage.Create(orderID, *data.Accrual)
				if err != nil {
					logger.Log.Error("Не удалось создать начисление", zap.Error(err))
					return
				}

				logger.Log.Info("Начисление сохранено",
					zap.String("orderID", orderID),
					zap.String("status", string(data.Status)),
					zap.Float64("accrual", *data.Accrual),
				)
			}

		default:
			logger.Log.Error("Статус не определен", zap.String("status", string(data.Status)))
		}
	})
}

func (as *AccrualService) StartCalculation() error {
	orders, err := as.Storage.FindAllUnprocessedOrders()
	if err != nil {
		return fmt.Errorf("ошибка при поиске необработанных заказов: %w", err)
	}

	if orders == nil {
		return nil
	}

	for _, order := range *orders {
		as.Calculate(order.ID)
	}

	return nil
}

func fetchAccrualData(endpoint string, orderID string) (data *accrualDataResponse, retryAfter int, err error) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/orders/%s", endpoint, orderID), nil)
	if err != nil {
		return nil, 0, fmt.Errorf("не удалось создать запрос: %w", err)
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("не удалось отправить GET-запрос: %w", err)
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusNoContent:
		return nil, 0, errNoOrder
	case http.StatusTooManyRequests:
		retryAfterHeader := res.Header.Get("Retry-After")
		retryAfter, err := strconv.Atoi(retryAfterHeader)
		if err != nil {
			retryAfter = defaultRetryAfterDuration
		}
		return nil, retryAfter, nil
	case http.StatusInternalServerError:
		return nil, 0, errServer
	}

	var parsedData accrualDataResponse
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(res.Body); err != nil {
		return nil, 0, fmt.Errorf("не удалось прочитать тело ответа: %w", err)
	}

	if err := json.Unmarshal(buf.Bytes(), &parsedData); err != nil {
		return nil, 0, fmt.Errorf("не удалось распаковать данные JSON: %w", err)
	}

	return &parsedData, 0, nil
}
