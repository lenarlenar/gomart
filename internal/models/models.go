package models

type UserRequest struct {
	Username string `json: "username"; binding:"required"`
	Password string `json: "username"; binding:"required"`
}

type User struct {
	Login string
	HashPass string
}