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

func TestLogin(t *testing.T) {
	router := setup.SetupRouter()
	setup.SetupTestData()

	// Contoh payload JSON untuk login
	payload := `{
		"phone_number": "0811111111",
		"pin": "123456"
	}`

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer([]byte(payload)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "SUCCESS", response["status"])
	assert.NotNil(t, response["result"])

	result := response["result"].(map[string]interface{})
	assert.NotEmpty(t, result["access_token"])
}
