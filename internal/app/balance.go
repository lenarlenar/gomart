package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lenarlenar/gomart/internal/models"
)

func (app *App) GetBalance(c *gin.Context) {
	userID := c.MustGet("user_id").(string)
	balance, err := app.BalanceService.GetUserBalance(userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, balance)
}

func (app *App) CreateWithdrawal(c *gin.Context) {
	var req models.Withdrawal
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный запрос"})
		return
	}

	if req.ID == nil || req.Sum == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "в запросе отсутствует номер заказа или сумма"})
		return
	}

	if *req.ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "номер заказа пустой"})
		return
	}

	if !app.OrdersService.Check(*req.ID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "номер заказа неверный"})
		return
	}

	userID := c.MustGet("user_id").(string)
	balance, err := app.BalanceService.GetUserBalance(userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if balance.Current < *req.Sum {
		c.JSON(http.StatusPaymentRequired, gin.H{"error": "недостаточно средств на балансе"})
		return
	}

	if err := app.BalanceService.CreateWithdrawal(*req.ID, userID, *req.Sum); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "успешно"})
}

func (app *App) GetWithdrawals(c *gin.Context) {
	userID := c.MustGet("user_id").(string)
	withdrawalFlow, err := app.BalanceService.GetWithdrawalFlow(userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(withdrawalFlow) == 0 {
		c.JSON(http.StatusNoContent, gin.H{"message": "списания не обнаружены"})
		return
	}

	c.JSON(http.StatusOK, withdrawalFlow)
}
