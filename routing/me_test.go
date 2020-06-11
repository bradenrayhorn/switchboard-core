package routing

import (
	"encoding/json"
	"fmt"
	"github.com/bradenrayhorn/switchboard-backend/repositories"
	"github.com/bradenrayhorn/switchboard-backend/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

type GetMeResponse struct {
	Id string `json:"id"`
}

func TestShowMe(t *testing.T) {
	r := MakeRouter()

	repositories.User = &repositories.MockUserRepository{}
	user, _ := repositories.User.CreateUser("test", "$2a$10$naqzJWUaOFm1/512Od.wPO4H8Vh8K38IGAb7rtgFizSflLVhpgMRG")

	token, _ := utils.CreateToken(user)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/me", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var body GetMeResponse
	_ = json.Unmarshal(w.Body.Bytes(), &body)
	assert.Equal(t, user.ID.Hex(), body.Id)
}

func TestCannotShowMeUnauthenticated(t *testing.T) {
	r := MakeRouter()

	repositories.User = &repositories.MockUserRepository{}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/me", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
