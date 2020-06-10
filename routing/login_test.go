package routing

import (
	"github.com/bradenrayhorn/switchboard-backend/repositories"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLogin(t *testing.T) {
	r := MakeRouter()

	repositories.User = &repositories.MockUserRepository{}
	_, _ = repositories.User.CreateUser("test", "$2a$10$naqzJWUaOFm1/512Od.wPO4H8Vh8K38IGAb7rtgFizSflLVhpgMRG")

	w := httptest.NewRecorder()
	reader := strings.NewReader("username=test&password=password")
	req, _ := http.NewRequest("POST", "/api/auth/login", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusOK)
}
