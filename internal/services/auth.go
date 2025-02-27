package services

import (
	"errors"
	"fmt"

	"github.com/lenarlenar/gomart/internal/interfaces"
	"github.com/lenarlenar/gomart/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	AuthStorage interfaces.AuthStorage
}


func (s *AuthService) Register(user models.UserRequest) error {

	hashPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if(err != nil) {
		return fmt.Errorf("ошибка при хеширований пароля: %w", err)
	}

	return s.AuthStorage.CreateUser(user.Username, string(hashPass))
}

func (s *AuthService) Login(userRequest models.UserRequest) error{
	user, err := s.AuthStorage.GetUser(userRequest.Username)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.HashPass), []byte(userRequest.Password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return models.ErrPasswordOrUsernameIsIncorrect
		}
		return fmt.Errorf("ошибка при сравнении паролей: %w", err)
	}

	return nil

}