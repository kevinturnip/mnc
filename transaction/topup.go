package transaction

import (
	"mnc/db"
	"mnc/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TopUp(c *gin.Context) {
	var topUpData struct {
		Amount int `json:"amount"`
	}
	if err := c.ShouldBindJSON(&topUpData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	phoneNumber := c.GetString("phone_number")
	var user model.User
	if err := db.DB.Where("phone_number = ?", phoneNumber).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User not found"})
		return
	}

	// Update balance
	balanceBefore := user.Balance
	user.Balance += topUpData.Amount
	db.DB.Save(&user)
	// topUpID := uuid.New().String()
	// Record the transaction
	transaction := model.Transaction{
		TransactionID: uuid.New().String(),
		// TopUpID: uuid.New().String(),
		UserID: user.UserID,
		// TopUpID:       topUpID,
		Type:          "CREDIT",
		Amount:        topUpData.Amount,
		BalanceBefore: balanceBefore,
		BalanceAfter:  user.Balance,
		CreatedDate:   time.Now(),
	}
	// type TopUp struct {
	// 	TopupID string `json:"top_up_id"`
	// 	*model.Transaction
	// }
	// var topup TopUp
	// b, _ := json.Marshal(transaction)
	// json.Unmarshal(b, &topup)
	// topup.TopupID = transaction.TransactionID
	db.DB.Create(&transaction)

	c.JSON(http.StatusOK, gin.H{
		"status": "SUCCESS",
		"result": transaction,
	})
}
