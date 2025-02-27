package app

import (
	"github.com/gin-gonic/gin"
	"github.com/lenarlenar/gomart/internal/interfaces"
)

type App struct {
	AuthStorage interfaces.AuthStorage
	AuthService interfaces.AuthService
	JWTService interfaces.JWTService
}

func (app *App) SetupRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Привет! ТЗ по данному проекту тут: https://github.com/lenarlenar/gomart",
		})
	})

	router.POST("api/user/register", app.Register)
	router.POST("api/user/login", app.Login)

	return router
}
