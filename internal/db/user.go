package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lenarlenar/gomart/internal/models"
	"github.com/lib/pq"
)

var (
	ErrDuplicateUser = errors.New("пользователь уже существует")
)

func (db *DataBase) CreateUser(name string, hashPass string) error {
	_, err := db.DB.ExecContext(context.Background(), "INSERT INTO users (login, hash) VALUES ($1, $2)", name, hashPass)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return models.ErrDuplicateUser
		}
		return fmt.Errorf("ошибка при попытке добавить пользователя в бд: %w", err)
	}
	return nil
}

func (db *DataBase) GetUser(username string) (*models.User, error) {
	user := new(models.User)
	err := db.DB.QueryRowContext(context.Background(), "SELECT id, login, hash FROM users WHERE login = $1", username).Scan(&user.ID, &user.Login, &user.HashPass)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrPasswordOrUsernameIsIncorrect
		}
		return nil, fmt.Errorf("ошибка при пойске пользователя в бд: %w", err)
	}

	return user, nil
}
