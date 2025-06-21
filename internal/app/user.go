package app

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lenarlenar/gomart/internal/models"
)

func (app *App) Register(c *gin.Context) {
	var userRequest models.UserRequest
	if err := c.ShouldBindJSON(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "невалидные входные значения",
		})
		return
	}

	if userRequest.Password == "" || userRequest.Login == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "поля username и password должны быть заполнены",
		})
		return
	}

	token, err := app.JWTService.Generate(userRequest.Login)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("ошибка при генерации токена %s", err.Error()),
		})
		return
	}

	err = app.AuthService.Register(userRequest)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateUser) {
			c.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("произошла ошибка при регистрации: %s", err.Error()),
			})
		}
		return
	}

	c.Header("Authorization", "Bearer "+token)
	c.JSON(http.StatusOK, gin.H{
		"message": "пользователь успешно зарегистрирован",
	})
}

func (app *App) Login(c *gin.Context) {
	var userRequest models.UserRequest
	if err := c.ShouldBindJSON(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "невалидные входные значения",
		})
		return
	}

	if userRequest.Password == "" || userRequest.Login == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "поля username и password должны быть заполнены",
		})
		return
	}

	if err := app.AuthService.Login(userRequest); err != nil {
		if errors.Is(err, models.ErrPasswordOrUsernameIsIncorrect) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}
		return
	}

	token, err := app.JWTService.Generate(userRequest.Login)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("ошибка при генерации jwt токена: %s", err.Error()),
		})
		return
	}

	c.Header("Authorization", "Bearer "+token)
	c.JSON(http.StatusOK, gin.H{
		"message": "вход выполнен",
	})
}
