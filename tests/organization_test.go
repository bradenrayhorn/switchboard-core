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
	_ = makeTestOrganizations(t, []*models.User{user1})

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

func TestAddUserToOrganization(t *testing.T) {
	user1, user2, token := makeTestUsersAndToken(t)
	organization := makeTestOrganizations(t, []*models.User{user1})

	w := httptest.NewRecorder()
	form := url.Values{"username": []string{user2.Username}, "organization_id": []string{organization.ID.Hex()}}
	req, _ := http.NewRequest("POST", "/api/organizations/invite-user", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	foundOrg, err := repositories.Organization.GetForUserAndID(organization.ID, user2.ID)
	assert.Nil(t, err)
	assert.NotNil(t, foundOrg)
	assert.Nil(t, repositories.User.DropAll())
	assert.Nil(t, repositories.Organization.DropAll())
}

func TestCannotAddUserToOrganizationWithEmptyParams(t *testing.T) {
	_, _, token := makeTestUsersAndToken(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/organizations/invite-user", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	assert.Nil(t, repositories.User.DropAll())
}

func TestCannotAddUserToInvalidOrganization(t *testing.T) {
	_, _, token := makeTestUsersAndToken(t)

	w := httptest.NewRecorder()
	form := url.Values{"username": []string{"test_user"}, "organization_id": []string{"invalid id"}}
	req, _ := http.NewRequest("POST", "/api/organizations/invite-user", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	assert.Nil(t, repositories.User.DropAll())
}

func TestCannotAddUserToOtherOrganization(t *testing.T) {
	user1, user2, token := makeTestUsersAndToken(t)
	organization := makeTestOrganizations(t, []*models.User{user2})

	w := httptest.NewRecorder()
	form := url.Values{"username": []string{user1.Username}, "organization_id": []string{organization.ID.Hex()}}
	req, _ := http.NewRequest("POST", "/api/organizations/invite-user", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	assert.Nil(t, repositories.User.DropAll())
	assert.Nil(t, repositories.Organization.DropAll())
}

func TestCannotAddUnknownUserToOrganization(t *testing.T) {
	user1, _, token := makeTestUsersAndToken(t)
	organization := makeTestOrganizations(t, []*models.User{user1})

	w := httptest.NewRecorder()
	form := url.Values{"username": []string{"weird username"}, "organization_id": []string{organization.ID.Hex()}}
	req, _ := http.NewRequest("POST", "/api/organizations/invite-user", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	assert.Nil(t, repositories.User.DropAll())
	assert.Nil(t, repositories.Organization.DropAll())
}

func TestCannotAddUserToOrganizationTwice(t *testing.T) {
	user1, _, token := makeTestUsersAndToken(t)
	organization := makeTestOrganizations(t, []*models.User{user1})

	w := httptest.NewRecorder()
	form := url.Values{"username": []string{user1.ID.Hex()}, "organization_id": []string{organization.ID.Hex()}}
	req, _ := http.NewRequest("POST", "/api/organizations/invite-user", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	assert.Nil(t, repositories.User.DropAll())
	assert.Nil(t, repositories.Organization.DropAll())
}

func makeTestOrganizations(t *testing.T, users []*models.User) *models.Organization {
	orgUsers := make([]models.OrganizationUser, 0)
	for _, user := range users {
		orgUsers = append(orgUsers, models.OrganizationUser{
			ID:       user.ID,
			Role:     models.RoleAdmin,
			Username: user.Username,
		})
	}
	organization, err := repositories.Organization.Create("Test Organization", orgUsers)
	assert.Nil(t, err)
	return organization
}
