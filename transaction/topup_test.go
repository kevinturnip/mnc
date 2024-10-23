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

func TestTopUp(t *testing.T) {
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

	// Top up request
	topUpPayload := `{"amount": 500000}`
	req, _ = http.NewRequest("POST", "/topup", bytes.NewBuffer([]byte(topUpPayload)))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "SUCCESS", response["status"])

	result := response["result"].(map[string]interface{})
	assert.Equal(t, float64(500000), result["amount_top_up"])
}
