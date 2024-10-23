package transaction

import (
	"mnc/db"
	"mnc/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Payment(c *gin.Context) {
	var paymentData struct {
		Amount  int    `json:"amount"`
		Remarks string `json:"remarks"`
	}
	if err := c.ShouldBindJSON(&paymentData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	phoneNumber := c.GetString("phone_number")
	var user model.User
	if err := db.DB.Where("phone_number = ?", phoneNumber).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User not found"})
		return
	}

	if user.Balance < paymentData.Amount {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Balance is not enough"})
		return
	}

	balanceBefore := user.Balance
	user.Balance -= paymentData.Amount
	db.DB.Save(&user)

	// Record the transaction
	transaction := model.Transaction{
		TransactionID: uuid.New().String(),
		UserID:        user.UserID,
		Type:          "DEBIT",
		Amount:        paymentData.Amount,
		Remarks:       paymentData.Remarks,
		BalanceBefore: balanceBefore,
		BalanceAfter:  user.Balance,
		CreatedDate:   time.Now(),
	}
	db.DB.Create(&transaction)

	c.JSON(http.StatusOK, gin.H{
		"status": "SUCCESS",
		"result": transaction,
	})
}
