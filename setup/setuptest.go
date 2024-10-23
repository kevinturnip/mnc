package setup

import (
	"mnc/auth"
	"mnc/db"
	jwttoken "mnc/jwt"
	"mnc/model"
	"mnc/profile"
	"mnc/transaction"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func SetupRouter() *gin.Engine {
	// Gunakan mode test untuk Gin
	gin.SetMode(gin.TestMode)

	// Inisialisasi database in-memory SQLite untuk testing
	db.DB, _ = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.DB.AutoMigrate(&model.User{}, &model.Transaction{})

	router := gin.Default()

	router.POST("/register", auth.Register)
	router.POST("/login", auth.Login)

	auth := router.Group("/")
	auth.Use(jwttoken.AuthenticateJWT())
	{
		auth.POST("/topup", transaction.TopUp)
		auth.POST("/pay", transaction.Payment)
		auth.POST("/transfer", transaction.Transfer)
		auth.GET("/transactions", transaction.TransactionReport)
		auth.PUT("/profile", profile.UpdateProfile)
	}

	return router
}

// Fungsi untuk setup data awal
func SetupTestData() {
	db.DB.Create(&model.User{
		UserID:      "e26c138d-c617-4949-a2b1-d1c283715a98",
		FirstName:   "petir",
		LastName:    "mandala",
		PhoneNumber: "0811111111",
		Pin:         "123456",
		Balance:     1000000,
	})

}
