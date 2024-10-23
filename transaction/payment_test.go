package transaction_test

import (
	"bytes"
	"encoding/json"
	"mnc/setup"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPayment(t *testing.T) {
	router := setup.SetupRouter()
	setup.SetupTestData()

	// Lakukan login untuk mendapatkan JWT token
	payload := `{"phone_number": "0811111111", "pin": "123456"}`
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer([]byte(payload)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var loginResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &loginResponse)
	token := loginResponse["result"].(map[string]interface{})["access_token"].(string)

	// Payment request
	paymentPayload := `{"amount": 100000, "remarks": "Payment test"}`
	req, _ = http.NewRequest("POST", "/pay", bytes.NewBuffer([]byte(paymentPayload)))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "SUCCESS", response["status"])
}
