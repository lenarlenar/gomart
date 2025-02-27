package router

import (
	"github.com/gin-gonic/gin"
	"github.com/lenarlenar/gomart/internal/handler"
)

func Run(addr string) {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Привет! ТЗ по данному проекту тут: https://github.com/lenarlenar/gomart",
		})
	})

	router.POST("api/user/register", handler.Register)
	router.POST("api/user/login", handler.Login)
	router.POST("api/user/logout", handler.Logout)

	router.Run(addr)
}
