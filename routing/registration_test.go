package routing

import (
	"github.com/bradenrayhorn/switchboard-core/repositories"
	"github.com/bradenrayhorn/switchboard-core/repositories/mocks"
	"github.com/bradenrayhorn/switchboard-core/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRegister(t *testing.T) {
	r := MakeTestRouter()
	utils.SetupTestRsaKeys()

	userRepo := new(mocks.UserRepository)
	userRepo.On("Exists", "test").Return(false, nil)
	userRepo.On("CreateUser", "test", mock.Anything).Return(utils.MakeTestUser("test", ""), nil)

	repositories.User = userRepo

	w := httptest.NewRecorder()
	reader := strings.NewReader("username=test&password=password")
	req, _ := http.NewRequest("POST", "/api/auth/register", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	userRepo.AssertCalled(t, "CreateUser", "test", mock.Anything)
}

func TestCannotRegisterTwice(t *testing.T) {
	r := MakeTestRouter()
	utils.SetupTestRsaKeys()
	userRepo := new(mocks.UserRepository)
	userRepo.On("Exists", "test").Return(true, nil)
	repositories.User = userRepo

	w := httptest.NewRecorder()
	reader := strings.NewReader("username=test&password=password")
	req, _ := http.NewRequest("POST", "/api/auth/register", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

	userRepo.AssertNotCalled(t, "CreateUser", "test", mock.Anything)
}

func TestCannotRegisterWithNoData(t *testing.T) {
	r := MakeTestRouter()
	utils.SetupTestRsaKeys()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/auth/register", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}
