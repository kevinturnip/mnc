package auth

import (
	"mnc/db"
	"mnc/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Register(c *gin.Context) {
	// Struct untuk menerima data dari request JSON
	var newUser struct {
		FirstName   string `json:"first_name" binding:"required"`
		LastName    string `json:"last_name" binding:"required"`
		PhoneNumber string `json:"phone_number" binding:"required"`
		Address     string `json:"address" binding:"required"`
		Pin         string `json:"pin" binding:"required"`
	}

	// Validasi input JSON
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	// Cek apakah nomor telepon sudah terdaftar
	var existingUser model.User
	if err := db.DB.Where("phone_number = ?", newUser.PhoneNumber).First(&existingUser).Error; err == nil {
		// Jika ada pengguna dengan nomor telepon tersebut
		c.JSON(http.StatusConflict, gin.H{"message": "Phone Number already registered"})
		return
	}

	// Buat pengguna baru
	user := model.User{
		UserID:      uuid.New().String(), // Generate UUID untuk user_id
		FirstName:   newUser.FirstName,
		LastName:    newUser.LastName,
		PhoneNumber: newUser.PhoneNumber,
		Address:     newUser.Address,
		Pin:         newUser.Pin,
		Balance:     0, // Balance awal
		CreatedDate: time.Now(),
	}

	// Simpan pengguna baru ke database
	if err := db.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not register user"})
		return
	}

	// Response berhasil
	c.JSON(http.StatusOK, gin.H{
		"status": "SUCCESS",
		"result": gin.H{
			"user_id":      user.UserID,
			"first_name":   user.FirstName,
			"last_name":    user.LastName,
			"phone_number": user.PhoneNumber,
			"address":      user.Address,
			"created_date": user.CreatedDate.Format("2006-01-02 15:04:05"),
		},
	})
}
