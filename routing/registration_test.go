package routing

import (
	"github.com/bradenrayhorn/switchboard-backend/repositories"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRegister(t *testing.T) {
	r := MakeRouter()

	repositories.User = &repositories.MockUserRepository{}

	w := httptest.NewRecorder()
	reader := strings.NewReader("username=test&password=password")
	req, _ := http.NewRequest("POST", "/api/auth/register", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	user, err := repositories.User.GetUser("test")

	assert.NotNil(t, user)
	assert.Nil(t, err)
}
