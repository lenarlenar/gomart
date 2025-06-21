package services

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	SecretKey string
}

func (s *JWTService) Generate(subject string) (string, error) {
	now := time.Now()
	exp := now.Add(time.Hour * 24).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": subject,
		"iat": now,
		"exp": exp,
	})

	signedToken, err := token.SignedString([]byte(s.SecretKey))

	if err != nil {
		return "", fmt.Errorf("ошибка при генерации jwt: %w", err)
	}

	return signedToken, nil
}

func (s *JWTService) Validate(tokenString string) (*jwt.Token, error) {
	parsedToken, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("не правильный метод подписи токена: %v", t.Header["alg"])
		}

		return []byte(s.SecretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("ошибка при валидации токена: %w", err)
	}

	if !parsedToken.Valid {
		return nil, jwt.ErrTokenInvalidId
	}

	return parsedToken, nil
}
