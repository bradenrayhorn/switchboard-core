package tests

import (
	"encoding/json"
	"fmt"
	"github.com/bradenrayhorn/switchboard-core/models"
	"github.com/bradenrayhorn/switchboard-core/repositories"
	"github.com/bradenrayhorn/switchboard-core/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestCreateGroup(t *testing.T) {
	user1, user2, token := makeTestUsersAndToken(t)
	organization := makeTestOrganizations(t, []*models.User{user1, user2})

	w := httptest.NewRecorder()
	form := url.Values{"users": []string{user1.ID.Hex(), user2.ID.Hex()}, "organization_id": []string{organization.ID.Hex()}}
	req, _ := http.NewRequest("POST", "/api/groups/create", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Nil(t, repositories.User.DropAll())
	assert.Nil(t, repositories.Group.DropAll())
	assert.Nil(t, repositories.Organization.DropAll())
}

func TestCantCreateDuplicateGroup(t *testing.T) {
	user1, user2, token := makeTestUsersAndToken(t)
	organization := makeTestOrganizations(t, []*models.User{user1, user2})

	_, err := repositories.Group.CreateGroup(nil, []primitive.ObjectID{user1.ID, user2.ID}, organization.ID)
	assert.Nil(t, err)

	w := httptest.NewRecorder()
	form := url.Values{"users": []string{user1.ID.Hex(), user2.ID.Hex()}, "organization_id": []string{organization.ID.Hex()}}
	req, _ := http.NewRequest("POST", "/api/groups/create", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	assert.Nil(t, repositories.User.DropAll())
	assert.Nil(t, repositories.Group.DropAll())
	assert.Nil(t, repositories.Organization.DropAll())
}

func TestCanDuplicateGroupOutsideOrganization(t *testing.T) {
	user1, user2, token := makeTestUsersAndToken(t)
	organization1 := makeTestOrganizations(t, []*models.User{user1, user2})
	organization2 := makeTestOrganizations(t, []*models.User{user1, user2})

	_, err := repositories.Group.CreateGroup(nil, []primitive.ObjectID{user1.ID, user2.ID}, organization1.ID)
	assert.Nil(t, err)

	w := httptest.NewRecorder()
	form := url.Values{"users": []string{user1.ID.Hex(), user2.ID.Hex()}, "organization_id": []string{organization2.ID.Hex()}}
	req, _ := http.NewRequest("POST", "/api/groups/create", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Nil(t, repositories.User.DropAll())
	assert.Nil(t, repositories.Group.DropAll())
	assert.Nil(t, repositories.Organization.DropAll())
}

func TestCantCreateGroupWithoutMe(t *testing.T) {
	_, user2, token := makeTestUsersAndToken(t)

	w := httptest.NewRecorder()
	form := url.Values{"users": []string{user2.ID.Hex()}}
	req, _ := http.NewRequest("POST", "/api/groups/create", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	assert.Nil(t, repositories.User.DropAll())
}

func TestCantCreateGroupWithoutUsers(t *testing.T) {
	_, _, token := makeTestUsersAndToken(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/groups/create", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	assert.Nil(t, repositories.User.DropAll())
}

func TestCantCreateGroupWithoutOrganization(t *testing.T) {
	_, user2, token := makeTestUsersAndToken(t)

	w := httptest.NewRecorder()
	form := url.Values{"users": []string{user2.ID.Hex()}}
	req, _ := http.NewRequest("POST", "/api/groups/create", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	assert.Nil(t, repositories.User.DropAll())
}

func TestGetGroups(t *testing.T) {
	user1, user2, token := makeTestUsersAndToken(t)
	organization := makeTestOrganizations(t, []*models.User{user1, user2})

	_, err := repositories.Group.CreateGroup(nil, []primitive.ObjectID{user1.ID, user2.ID}, organization.ID)
	assert.Nil(t, err)
	_, err = repositories.Group.CreateGroup(nil, []primitive.ObjectID{user1.ID}, organization.ID)
	assert.Nil(t, err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/groups", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var body struct {
		Data []interface{}
	}
	_ = json.Unmarshal(w.Body.Bytes(), &body)
	assert.Len(t, body.Data, 2)
}

func TestUpdateGroup(t *testing.T) {
	user1, user2, token := makeTestUsersAndToken(t)
	user3, err := repositories.User.CreateUser("test3", "")
	organization := makeTestOrganizations(t, []*models.User{user1, user2, user3})
	assert.Nil(t, err)

	group, err := repositories.Group.CreateGroup(nil, []primitive.ObjectID{user1.ID, user2.ID}, organization.ID)
	assert.Nil(t, err)

	w := httptest.NewRecorder()
	form := url.Values{
		"id":           []string{group.ID.Hex()},
		"name":         []string{"new group name"},
		"add_users":    []string{user3.ID.Hex()},
		"remove_users": []string{user2.ID.Hex()},
	}
	req, _ := http.NewRequest("POST", "/api/groups/update", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Nil(t, repositories.User.DropAll())
	assert.Nil(t, repositories.Group.DropAll())
}

func makeTestUsersAndToken(t *testing.T) (*models.User, *models.User, string) {
	user1, err := repositories.User.CreateUser("test1", "")
	assert.Nil(t, err)
	user2, err := repositories.User.CreateUser("test2", "")
	assert.Nil(t, err)
	token, _ := utils.CreateToken(user1)
	return user1, user2, token
}
