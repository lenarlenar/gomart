package interfaces

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lenarlenar/gomart/internal/models"
)

//go:generate mockgen -destination=mocks/mock_auth.go . AuthService
type AuthService interface {
	Register(user models.UserRequest) error
	Login(user models.UserRequest) error
}

//go:generate mockgen -destination=mocks/mock_auth_storage.go . AuthStorage
type AuthStorage interface {
	CreateUser(name string, hashPass string) error
	GetUser(username string) (*models.User, error)
}

//go:generate mockgen -destination=mocks/mock_jwt.go . JWTService
type JWTService interface {
	Generate(subject string) (string, error)
	Validate(token string) (*jwt.Token, error)
}

//go:generate mockgen -destination=mocks/mock_order.go . OrdersService
type OrdersService interface {
	Check(orderID string) bool
	CreateOrder(orderID, userID string) error
	GetOrders(userID string) ([]models.Order, error)
}

type OrdersStorage interface {
	CreateOrder(orderID, userID string) error
	FindOrder(orderID string) (*models.OrderDB, error)
	FindOrdersWithAccrual(userID string) (*[]models.OrderWithAccrualDB, error)
}

//go:generate mockgen -destination=mocks/mock_accrual.go . AccrualService
type AccrualService interface {
	Calculate(orderID string)
	StartCalculation() error
}

type AccrualStorage interface {
	UpdateStatus(orderID string, status models.OrderStatusDB) error
	Create(orderID string, amount float64) error
	FindAllUnprocessedOrders() (*[]models.OrderDB, error)
}

type AccrualJobQueueService interface {
	Enqueue(job models.Job)
	ScheduleJob(job models.Job, delay time.Duration)
	PauseAndResume(delay time.Duration)
}

type BalanceService interface {
	GetUserBalance(userID string) (models.Balance, error)
	CreateWithdrawal(orderID, userID string, amount float64) error
	GetWithdrawalFlow(userID string) ([]models.WithdrawalFlowItem, error)
}

type BalanceStorage interface {
	FindAccrualFlow(userID string) (*[]models.AccrualFlowItemDB, error)
	CreateWithdrawal(orderID, userID string, amount float64) error
	FindWithdrawalFlow(userID string) (*[]models.WithdrawalFlowItemDB, error)
}
