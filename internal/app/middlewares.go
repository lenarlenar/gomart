package app

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (app *App) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/api/user/register" || c.Request.URL.Path == "/api/user/login" {
			c.Next()
			return
		}

		authHeader := c.GetHeader("Authorization")
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if authHeader == "" || tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization заголовок отсутствует",
			})
			c.Abort()
			return
		}

		token, err := app.JWTService.Validate(tokenString)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": fmt.Sprintf("ошибка при проверке токена: %s", err.Error()),
			})
			c.Abort()
			return
		}

		login, err := token.Claims.GetSubject()
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": fmt.Sprintf("ошибка при проверке токена: %s", err.Error()),
			})
			c.Abort()
			return
		}

		user, err := app.AuthStorage.GetUser(login)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": fmt.Sprintf("пользователь не аутентифицирован: %s", err.Error()),
			})
			c.Abort()
			return
		}

		c.Set("user_name", user.Login)
		c.Set("user_id", user.ID)
		c.Next()
	}
}
