package transaction

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mnc/db"
	"mnc/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var rdb *redis.Client
var ctx = context.Background()

// Inisialisasi Redis client
func InitRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Ganti dengan alamat Redis Anda
	})
}

// Struct untuk menyimpan data transfer
type TransferRequest struct {
	SourcePhoneNumber string
	TargetUserID      string
	Amount            int
	Remarks           string
}

// Fungsi untuk menambahkan permintaan transfer ke antrean Redis
func EnqueueTransfer(request TransferRequest) error {
	data, err := json.Marshal(request)
	if err != nil {
		return err
	}
	log.Printf("Enqueuing transfer: %+v\n", request)
	return rdb.LPush(ctx, "transfer_queue", data).Err()
}

// Worker untuk memproses antrean transfer
func TransferWorker() {
	for {
		data, err := rdb.BRPop(ctx, 0*time.Second, "transfer_queue").Result()
		if err != nil {
			log.Println("Error reading from Redis queue:", err)
			continue
		}

		var request TransferRequest
		if err := json.Unmarshal([]byte(data[1]), &request); err != nil {
			log.Println("Error unmarshalling transfer request:", err)
			continue
		}
		log.Printf("Processing transfer: %+v\n", request)

		if err := processTransfer(request); err != nil {
			log.Println("Error processing transfer:", err)
		} else {
			log.Println("Transfer processed successfully:", request)
		}
	}
}

// Fungsi untuk memproses transfer
func processTransfer(request TransferRequest) error {
	log.Println("Fetching users for transfer...") // Log tahap awal

	var sourceUser, targetUser model.User
	if err := db.DB.Where("phone_number = ?", request.SourcePhoneNumber).First(&sourceUser).Error; err != nil {
		return fmt.Errorf("source user not found")
	}

	if err := db.DB.Where("user_id = ?", request.TargetUserID).First(&targetUser).Error; err != nil {
		return fmt.Errorf("target user not found")
	}

	if sourceUser.Balance < request.Amount {
		return fmt.Errorf("balance is not enough")
	}

	balanceBeforeSource := sourceUser.Balance
	balanceBeforeTarget := targetUser.Balance

	log.Printf("Balances before transfer: Source: %d, Target: %d\n", balanceBeforeSource, balanceBeforeTarget)

	return db.DB.Transaction(func(tx *gorm.DB) error {
		sourceUser.Balance -= request.Amount
		targetUser.Balance += request.Amount

		if err := tx.Save(&sourceUser).Error; err != nil {
			return fmt.Errorf("could not update source user balance: %v", err)
		}
		if err := tx.Save(&targetUser).Error; err != nil {
			return fmt.Errorf("could not update target user balance: %v", err)
		}

		sourceTransaction := model.Transaction{
			TransactionID: uuid.New().String(),
			UserID:        sourceUser.UserID,
			Type:          "DEBIT",
			Amount:        request.Amount,
			Remarks:       request.Remarks,
			BalanceBefore: balanceBeforeSource,
			BalanceAfter:  sourceUser.Balance,
			CreatedDate:   time.Now(),
		}
		if err := tx.Create(&sourceTransaction).Error; err != nil {
			return fmt.Errorf("could not create source transaction: %v", err)
		}

		targetTransaction := model.Transaction{
			TransactionID: uuid.New().String(),
			UserID:        targetUser.UserID,
			Type:          "CREDIT",
			Amount:        request.Amount,
			Remarks:       request.Remarks,
			BalanceBefore: balanceBeforeTarget,
			BalanceAfter:  targetUser.Balance,
			CreatedDate:   time.Now(),
		}
		if err := tx.Create(&targetTransaction).Error; err != nil {
			return fmt.Errorf("could not create target transaction: %v", err)
		}

		log.Println("Transfer successfully completed") // Log setelah transfer selesai
		return nil
	})
}

// Endpoint transfer
func Transfer(c *gin.Context) {
	var transferData struct {
		TargetUserID string `json:"target_user"`
		Amount       int    `json:"amount"`
		Remarks      string `json:"remarks"`
	}
	if err := c.ShouldBindJSON(&transferData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	phoneNumber := c.GetString("phone_number")
	var user model.User
	if err := db.DB.Where("phone_number = ?", phoneNumber).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User not found"})
		return
	}

	if transferData.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid transfer amount"})
		return
	}

	// Buat permintaan transfer
	request := TransferRequest{
		SourcePhoneNumber: phoneNumber,
		TargetUserID:      transferData.TargetUserID,
		Amount:            transferData.Amount,
		Remarks:           transferData.Remarks,
	}

	// Tambahkan permintaan transfer ke antrean Redis
	if err := EnqueueTransfer(request); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not enqueue transfer"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "SUCCESS",
		"result": "Transfer request has been enqueued",
	})
}
