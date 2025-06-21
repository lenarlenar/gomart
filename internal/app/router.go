package app

import (
	"github.com/gin-gonic/gin"
	"github.com/lenarlenar/gomart/internal/interfaces"
)

type App struct {
	AuthStorage    interfaces.AuthStorage
	AuthService    interfaces.AuthService
	JWTService     interfaces.JWTService
	OrdersService  interfaces.OrdersService
	AccrualService interfaces.AccrualService
	BalanceService interfaces.BalanceService
}

func (app *App) SetupRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Привет! ТЗ по данному проекту тут: https://github.com/lenarlenar/gomart",
		})
	})

	router.Use(app.AuthMiddleware())
	router.POST("api/user/register", app.Register)
	router.POST("api/user/login", app.Login)
	router.POST("api/user/orders", app.Orders)
	router.GET("api/user/orders", app.GetOrders)
	router.GET("api/user/balance", app.GetBalance)
	router.POST("api/user/balance/withdraw", app.CreateWithdrawal)
	router.GET("api/user/withdrawals", app.GetWithdrawals)

	return router
}
