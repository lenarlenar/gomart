package app

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lenarlenar/gomart/internal/models"
)

func (app *App) Orders(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusBadRequest, "ошибка чтения тела запроса")
		return
	}
	orderID := string(body)
	if orderID == "" {
		c.String(http.StatusUnprocessableEntity, "id заказа не может быть пустым")
		return
	}
	if !app.OrdersService.Check(orderID) {
		c.String(http.StatusUnprocessableEntity, "неверный формат номера заказа")
		return
	}

	userID := c.MustGet("user_id").(string)
	errCreateOrder := app.OrdersService.CreateOrder(orderID, userID)
	if errCreateOrder != nil {
		if errors.Is(errCreateOrder, models.ErrDuplicateOrderByUser) {
			c.String(http.StatusOK, "заказ уже был добавлен этим пользователем")
			return
		}
		if errors.Is(errCreateOrder, models.ErrDuplicateOrder) {
			c.String(http.StatusConflict, "заказ уже был создан другим пользователем")
			return
		}

		c.String(http.StatusInternalServerError, fmt.Sprintf("произошла ошибка при создании заказа: %s", errCreateOrder.Error()))
		return
	}

	app.AccrualService.Calculate(orderID)
	c.String(http.StatusAccepted, "новый номер заказа принят в обработку")
}

func (app *App) GetOrders(c *gin.Context) {
	userID := c.MustGet("user_id").(string)
	orders, err := app.OrdersService.GetOrders(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "произошла ошибка при получении заказов",
		})
		return
	}

	if len(orders) == 0 {
		c.JSON(http.StatusNoContent, gin.H{
			"message": "у пользователя нет заказов",
		})
		return
	}
	c.JSON(http.StatusOK, orders)
}
