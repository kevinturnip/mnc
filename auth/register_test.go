package auth_test

import (
	"bytes"
	"encoding/json"
	"mnc/setup"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	router := setup.SetupRouter()

	// Contoh payload JSON untuk register
	payload := `{
		"first_name": "Jane",
		"last_name": "Doe",
		"phone_number": "0811111111",
		"address": "Some Street",
		"pin": "123456"
	}`

	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer([]byte(payload)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "SUCCESS", response["status"])
	assert.NotNil(t, response["result"])
}
