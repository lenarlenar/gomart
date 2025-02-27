package models

import "errors"

var (
	ErrDuplicateUser                  = errors.New("пользователь уже существует")
	ErrPasswordOrUsernameIsIncorrect  = errors.New("неверный логин или пароль")
)