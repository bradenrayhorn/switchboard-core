package routing

import (
	"errors"
	"fmt"
	"github.com/bradenrayhorn/switchboard-core/repositories"
	"github.com/bradenrayhorn/switchboard-core/repositories/mocks"
	"github.com/bradenrayhorn/switchboard-core/utils"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLogin(t *testing.T) {
	userRepo := new(mocks.UserRepository)
	user := utils.MakeTestUser("test", "$2a$10$naqzJWUaOFm1/512Od.wPO4H8Vh8K38IGAb7rtgFizSflLVhpgMRG")
	userRepo.On("GetUser", "test").Return(user, nil)
	repositories.User = userRepo

	testLogin(t, http.StatusOK, "test", "password")
}

func TestCannotLoginWithInvalidUsername(t *testing.T) {
	userRepo := new(mocks.UserRepository)
	userRepo.On("GetUser", "test-bad").Return(nil, errors.New("user not found"))
	repositories.User = userRepo

	testLogin(t, http.StatusUnprocessableEntity, "test-bad", "password")
}

func TestCannotLoginWithInvalidPassword(t *testing.T) {
	userRepo := new(mocks.UserRepository)
	user := utils.MakeTestUser("test", "$2a$10$naqzJWUaOFm1/512Od.wPO4H8Vh8K38IGAb7rtgFizSflLVhpgMRG")
	userRepo.On("GetUser", "test").Return(user, nil)
	repositories.User = userRepo

	testLogin(t, http.StatusUnprocessableEntity, "test", "password-bad")
}

func testLogin(t *testing.T, expectedStatus int, username string, password string) {
	utils.SetupTestRsaKeys()
	r := MakeTestRouter()
	w := httptest.NewRecorder()
	reader := strings.NewReader(fmt.Sprintf("username=%s&password=%s", username, password))
	req, _ := http.NewRequest("POST", "/api/auth/login", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, req)

	assert.Equal(t, expectedStatus, w.Code)
}
