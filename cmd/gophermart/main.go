package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/lenarlenar/gomart/internal/app"
	"github.com/lenarlenar/gomart/internal/config"
	"github.com/lenarlenar/gomart/internal/db"
	"github.com/lenarlenar/gomart/internal/logger"
	"github.com/lenarlenar/gomart/internal/services"
)

func main() {

	if err := logger.Initialize("debug", "development"); err != nil {
		log.Fatalf("Ошибка инициализации логгера: %s", err)
	}

	config := config.Get()
	storage, err := db.Open(config.DSN)
	if err != nil {
		log.Fatalf("Ошибка при открытии подключения к бд: %v", err)
	}

	defer storage.DB.Close()

	if err := storage.RunMigrations(); err != nil {
		log.Fatalf("Ошибка при запуске миграций базы данных: %s", err)
	}

	jwt := &services.JWTService{SecretKey: "secret_key"}
	auth := &services.AuthService{AuthStorage: storage}
	orders := &services.OrdersService{OrdersStorage: storage}
	jobService := services.NewJobQueueService(context.Background(), 100, 5)
	handleTerminationProcess(func() {
		jobService.Shutdown()
	})
	
	accrualService := &services.AccrualService{Storage: storage, JobQueueService: jobService, ExternalEndpoint: config.AccrualAddress}

	if err := accrualService.StartCalculation(); err != nil {
		log.Fatalf("Ошибка при запуске расчета начислений: %s", err)
	}

	balanceService := &services.BalanceService{Storage: storage}

	app := app.App{
		AuthService:    auth,
		AuthStorage:    storage,
		JWTService:     jwt,
		OrdersService:  orders,
		AccrualService: accrualService,
		BalanceService: balanceService,
	}

	app.SetupRouter().Run(config.Address)
}

func handleTerminationProcess(cleanup func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cleanup()
		os.Exit(1)
	}()
}
