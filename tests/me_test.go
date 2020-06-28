package tests

import (
	"encoding/json"
	"fmt"
	"github.com/bradenrayhorn/switchboard-core/repositories"
	"github.com/bradenrayhorn/switchboard-core/utils"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type GetMeResponse struct {
	Id string `json:"id"`
}

func TestShowMe(t *testing.T) {
	user, err := repositories.User.CreateUser("test", "")
	assert.Nil(t, err)

	token, _ := utils.CreateToken(user)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/me", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var body GetMeResponse
	_ = json.Unmarshal(w.Body.Bytes(), &body)
	assert.Equal(t, user.ID.Hex(), body.Id)
	assert.Nil(t, repositories.User.DropAll())
}

func TestCannotShowMeUnauthenticated(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/me", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestCannotShowMeWithExpiredToken(t *testing.T) {
	user, err := repositories.User.CreateUser("test", "")
	assert.Nil(t, err)

	viper.Set("token_expiration", -10*time.Second)
	token, _ := utils.CreateToken(user)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/me", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Nil(t, repositories.User.DropAll())
	viper.Set("token_expiration", 24*time.Hour)
}
