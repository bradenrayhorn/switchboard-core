package routing

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
	repositories.User = &repositories.MockUserRepository{}
	_, _ = repositories.User.CreateUser("test", "$2a$10$naqzJWUaOFm1/512Od.wPO4H8Vh8K38IGAb7rtgFizSflLVhpgMRG")

	testLogin(t, http.StatusOK, "test", "password")
}

func TestCannotLoginWithInvalidUsername(t *testing.T) {
	repositories.User = &repositories.MockUserRepository{}
	_, _ = repositories.User.CreateUser("test", "$2a$10$naqzJWUaOFm1/512Od.wPO4H8Vh8K38IGAb7rtgFizSflLVhpgMRG")

	testLogin(t, http.StatusUnprocessableEntity, "test-bad", "password")
}

func TestCannotLoginWithInvalidPassword(t *testing.T) {
	repositories.User = &repositories.MockUserRepository{}
	_, _ = repositories.User.CreateUser("test", "$2a$10$naqzJWUaOFm1/512Od.wPO4H8Vh8K38IGAb7rtgFizSflLVhpgMRG")

	testLogin(t, http.StatusUnprocessableEntity, "test", "password-bad")
}

func testLogin(t *testing.T, expectedStatus int, username string, password string) {
	r := MakeTestRouter()
	w := httptest.NewRecorder()
	reader := strings.NewReader(fmt.Sprintf("username=%s&password=%s", username, password))
	req, _ := http.NewRequest("POST", "/api/auth/login", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, req)

	assert.Equal(t, expectedStatus, w.Code)
}
