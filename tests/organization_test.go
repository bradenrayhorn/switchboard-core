package tests

import (
	"encoding/json"
	"fmt"
	"github.com/bradenrayhorn/switchboard-core/models"
	"github.com/bradenrayhorn/switchboard-core/repositories"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestCreateOrganization(t *testing.T) {
	_, _, token := makeTestUsersAndToken(t)

	w := httptest.NewRecorder()
	form := url.Values{"name": []string{"test organization"}}
	req, _ := http.NewRequest("POST", "/api/organizations", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Nil(t, repositories.User.DropAll())
	assert.Nil(t, repositories.Organization.DropAll())
}

func TestCannotCreateOrganizationWithoutName(t *testing.T) {
	_, _, token := makeTestUsersAndToken(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/organizations", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	assert.Nil(t, repositories.User.DropAll())
}

func TestGetOrganizations(t *testing.T) {
	user1, _, token := makeTestUsersAndToken(t)
	_ = makeTestOrganizations(t, user1)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/organizations", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var body struct {
		Data []interface{}
	}
	_ = json.Unmarshal(w.Body.Bytes(), &body)
	assert.Len(t, body.Data, 1)
	assert.Nil(t, repositories.User.DropAll())
	assert.Nil(t, repositories.Organization.DropAll())
}

func makeTestOrganizations(t *testing.T, user *models.User) *models.Organization {
	users := []models.OrganizationUser{{
		ID:   user.ID,
		Role: models.RoleAdmin,
	}}
	organization, err := repositories.Organization.Create("Test Organization", users)
	assert.Nil(t, err)
	return organization
}
