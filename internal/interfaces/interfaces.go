package interfaces

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/lenarlenar/gomart/internal/models"
)

//go:generate mockgen -destination=mocks/mock_auth.go . AuthService
type AuthService interface {
	Register(user models.UserRequest) error
	Login(user models.UserRequest) error
}

//go:generate mockgen -destination=mocks/mock_storage.go . AuthStorage
type AuthStorage interface {
	CreateUser(name string, hashPass string) error
	GetUser(username string) (*models.User, error)
}

//go:generate mockgen -destination=mocks/mock_jwt.go . JWTService
type JWTService interface {
	Generate(subject string) (string, error)
	Validate(token string) (*jwt.Token, error)
}