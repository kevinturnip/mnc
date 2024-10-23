package transaction_test

import (
	"bytes"
	"encoding/json"
	"mnc/db"
	"mnc/model"
	"mnc/setup"
	"mnc/transaction"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransfer(t *testing.T) {
	// db.InitDB()
	transaction.InitRedis()
	go transaction.TransferWorker()
	router := setup.SetupRouter()
	setup.SetupTestData()

	// Buat pengguna tujuan untuk menerima transfer
	targetUser := model.User{
		UserID:      "e26c138d-c617-4949-a2b1-d1c283715a99",
		FirstName:   "hujan",
		LastName:    "badai",
		PhoneNumber: "0822222222",
		Pin:         "123456",
		Balance:     1000000,
	}
	db.DB.Create(&targetUser)

	// Lakukan login untuk mendapatkan JWT token
	payload := `{"phone_number": "0811111111", "pin": "123456"}`
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer([]byte(payload)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var loginResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &loginResponse)
	token := loginResponse["result"].(map[string]interface{})["access_token"].(string)

	// Transfer request
	transferPayload := `{
		"target_user": "e26c138d-c617-4949-a2b1-d1c283715a99",
		"amount": 200000,
		"remarks": "Test Transfer"
	}`
	req, _ = http.NewRequest("POST", "/transfer", bytes.NewBuffer([]byte(transferPayload)))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Pastikan status code adalah 200
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "SUCCESS", response["status"])

	// Cek apakah transfer berhasil dan saldo sudah terupdate
	var sourceUser, updatedTargetUser model.User
	db.DB.Where("phone_number = ?", "0811111111").First(&sourceUser)
	db.DB.Where("phone_number = ?", "0822222222").First(&updatedTargetUser)

	assert.Equal(t, 800000, sourceUser.Balance)         // Saldo pengguna sumber setelah transfer
	assert.Equal(t, 1200000, updatedTargetUser.Balance) // Saldo pengguna tujuan setelah transfer
}
