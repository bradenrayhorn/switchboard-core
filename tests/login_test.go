package tests

import (
	"fmt"
	"github.com/bradenrayhorn/switchboard-core/repositories"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLogin(t *testing.T) {
	_, err := repositories.User.CreateUser("test", "$2a$10$naqzJWUaOFm1/512Od.wPO4H8Vh8K38IGAb7rtgFizSflLVhpgMRG")
	assert.Nil(t, err)

	testLogin(t, http.StatusOK, "test", "password")

	assert.Nil(t, repositories.User.DropAll())
}

func TestCannotLoginWithInvalidUsername(t *testing.T) {
	testLogin(t, http.StatusUnprocessableEntity, "test-bad", "password")
}

func TestCannotLoginWithInvalidPassword(t *testing.T) {
	_, err := repositories.User.CreateUser("test", "$2a$10$naqzJWUaOFm1/512Od.wPO4H8Vh8K38IGAb7rtgFizSflLVhpgMRG")
	assert.Nil(t, err)

	testLogin(t, http.StatusUnprocessableEntity, "test", "password-wrong")

	assert.Nil(t, repositories.User.DropAll())
}

func testLogin(t *testing.T, expectedStatus int, username string, password string) {
	w := httptest.NewRecorder()
	reader := strings.NewReader(fmt.Sprintf("username=%s&password=%s", username, password))
	req, _ := http.NewRequest("POST", "/api/auth/login", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, req)

	assert.Equal(t, expectedStatus, w.Code)
}
