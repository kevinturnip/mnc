package profile

import (
	"mnc/db"
	"mnc/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func UpdateProfile(c *gin.Context) {
	var profileData struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Address   string `json:"address"`
	}
	if err := c.ShouldBindJSON(&profileData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	phoneNumber := c.GetString("phone_number")
	var user model.User
	if err := db.DB.Where("phone_number = ?", phoneNumber).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User not found"})
		return
	}

	user.FirstName = profileData.FirstName
	user.LastName = profileData.LastName
	user.Address = profileData.Address
	user.UpdatedDate = time.Now()
	db.DB.Save(&user)

	c.JSON(http.StatusOK, gin.H{
		"status": "SUCCESS",
		"result": user,
	})
}
