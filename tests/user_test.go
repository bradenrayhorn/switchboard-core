package tests

import (
	"encoding/json"
	"fmt"
	"github.com/bradenrayhorn/switchboard-core/models"
	"github.com/bradenrayhorn/switchboard-core/repositories"
	"github.com/bradenrayhorn/switchboard-core/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestUserSearch(t *testing.T) {
	george, err := repositories.User.CreateUser("george", "")
	assert.Nil(t, err)
	thomas, err := repositories.User.CreateUser("thomas", "")
	assert.Nil(t, err)
	token, _ := utils.CreateToken(george)
	organization := makeTestOrganizations(t, []*models.User{george, thomas})

	w := httptest.NewRecorder()
	form := url.Values{"username": []string{"thom"}, "organization_id": []string{organization.ID.Hex()}}
	req, _ := http.NewRequest("POST", "/api/users/search", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var body struct {
		Data []models.UserResponse
	}
	_ = json.Unmarshal(w.Body.Bytes(), &body)
	require.Len(t, body.Data, 1)
	assert.Equal(t, thomas.ID.Hex(), body.Data[0].ID.Hex())

	assert.Nil(t, repositories.User.DropAll())
}
