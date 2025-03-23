package models

import "errors"

var (
	ErrDuplicateUser                 = errors.New("пользователь уже существует")
	ErrPasswordOrUsernameIsIncorrect = errors.New("неверный логин или пароль")
	ErrDuplicateOrder                = errors.New("заказ уже существует")
	ErrDuplicateOrderByUser          = errors.New("заказ уже был создан этим пользователем")
)
