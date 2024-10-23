package transaction

import (
	"mnc/db"
	"mnc/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func TransactionReport(c *gin.Context) {
	phoneNumber := c.GetString("phone_number")
	var user model.User
	if err := db.DB.Where("phone_number = ?", phoneNumber).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User not found"})
		return
	}

	var transactions []model.Transaction
	db.DB.Where("user_id = ?", user.UserID).Find(&transactions)

	c.JSON(http.StatusOK, gin.H{
		"status": "SUCCESS",
		"result": transactions,
	})
}
