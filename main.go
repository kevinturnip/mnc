package main

import (
	"mnc/auth"
	"mnc/db"
	jwttoken "mnc/jwt"
	"mnc/profile"
	"mnc/transaction"

	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDB()
	transaction.InitRedis()
	go transaction.TransferWorker()
	router := gin.Default()

	router.POST("/register", auth.Register)
	router.POST("/login", auth.Login)

	auth := router.Group("/")
	auth.Use(jwttoken.AuthenticateJWT())
	{
		auth.POST("/topup", transaction.TopUp)
		auth.POST("/pay", transaction.Payment)
		auth.POST("/transfer", transaction.Transfer) // queue implementation
		auth.GET("/transactions", transaction.TransactionReport)
		auth.PUT("/profile", profile.UpdateProfile)
	}

	router.Run(":8080")
}
