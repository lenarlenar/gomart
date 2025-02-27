package main

import (
	"log"

	"github.com/lenarlenar/gomart/internal/app"
	"github.com/lenarlenar/gomart/internal/config"
	"github.com/lenarlenar/gomart/internal/db"
	"github.com/lenarlenar/gomart/internal/db/migrations"
	"github.com/lenarlenar/gomart/internal/services"
)

func main() {
	config := config.Get()
	storage, err := db.Open(config.DBUri)
	if(err != nil) {
		log.Fatalf("Ошибка при открытии подключения к бд: %v", err)
	}

	defer storage.DB.Close()
	migrations.Apply(storage.DB)
	jwt := &services.JWTService {SecretKey: "secret_key"}
	auth := &services.AuthService{AuthStorage: storage}
	app := app.App {
		AuthService: auth,
		AuthStorage: storage,
		JWTService: jwt,
	}
	app.SetupRouter().Run(config.Address)
}
