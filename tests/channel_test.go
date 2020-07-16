package tests

import (
	"bytes"
	"fmt"
	"github.com/bradenrayhorn/switchboard-core/models"
	"github.com/bradenrayhorn/switchboard-core/repositories"
	"github.com/bradenrayhorn/switchboard-core/utils"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateChannel(t *testing.T) {
	user1, user2, token := makeTestUsersAndToken(t)
	organization := makeTestOrganizations(t, []*models.User{user1, user2})

	w := httptest.NewRecorder()
	json := []byte(`{"name":"test-channel","organization_id":"` + organization.ID.Hex() + `","private":false}`)
	req, _ := http.NewRequest("POST", "/api/channels", bytes.NewBuffer(json))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Nil(t, repositories.User.DropAll())
	assert.Nil(t, repositories.Group.DropAll())
	assert.Nil(t, repositories.Organization.DropAll())
}

func TestCannotCreateChannelOutsideOrganization(t *testing.T) {
	_, user2, token := makeTestUsersAndToken(t)
	organization := makeTestOrganizations(t, []*models.User{user2})

	w := httptest.NewRecorder()
	json := []byte(`{"name":"test-channel","organization_id:` + organization.ID.Hex() + `","private":false}`)
	req, _ := http.NewRequest("POST", "/api/channels", bytes.NewBuffer(json))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	assert.Nil(t, repositories.User.DropAll())
	assert.Nil(t, repositories.Group.DropAll())
	assert.Nil(t, repositories.Organization.DropAll())
}

func TestCannotCreateChannelWithInvalidName(t *testing.T) {
	user1, user2, token := makeTestUsersAndToken(t)
	organization := makeTestOrganizations(t, []*models.User{user1, user2})

	w := httptest.NewRecorder()
	json := []byte(`{"name":"test channel","organization_id:` + organization.ID.Hex() + `","private":false}`)
	req, _ := http.NewRequest("POST", "/api/channels", bytes.NewBuffer(json))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	assert.Nil(t, repositories.User.DropAll())
	assert.Nil(t, repositories.Group.DropAll())
	assert.Nil(t, repositories.Organization.DropAll())
}

func TestCannotCreateChannelSameName(t *testing.T) {
	user1, user2, token := makeTestUsersAndToken(t)
	organization := makeTestOrganizations(t, []*models.User{user1, user2})

	makeTestChannel(organization.ID, false)

	w := httptest.NewRecorder()
	json := []byte(`{"name":"test-channel","organization_id:` + organization.ID.Hex() + `","private":false}`)
	req, _ := http.NewRequest("POST", "/api/channels", bytes.NewBuffer(json))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	assert.Nil(t, repositories.User.DropAll())
	assert.Nil(t, repositories.Group.DropAll())
	assert.Nil(t, repositories.Organization.DropAll())
}

func TestCanJoinChannel(t *testing.T) {
	user1, user2, token := makeTestUsersAndToken(t)
	organization := makeTestOrganizations(t, []*models.User{user1, user2})

	channel := makeTestChannel(organization.ID, false)

	w := httptest.NewRecorder()
	json := []byte(`{"channel_id":"` + channel.ID.Hex() + `"}`)
	req, _ := http.NewRequest("POST", "/api/channels/join", bytes.NewBuffer(json))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Nil(t, repositories.User.DropAll())
	assert.Nil(t, repositories.Group.DropAll())
	assert.Nil(t, repositories.Organization.DropAll())
}

func TestCannotJoinPrivateChannel(t *testing.T) {
	user1, user2, token := makeTestUsersAndToken(t)
	organization := makeTestOrganizations(t, []*models.User{user1, user2})

	channel := makeTestChannel(organization.ID, true)

	w := httptest.NewRecorder()
	json := []byte(`{"channel_id":"` + channel.ID.Hex() + `"}`)
	req, _ := http.NewRequest("POST", "/api/channels/join", bytes.NewBuffer(json))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	assert.Nil(t, repositories.User.DropAll())
	assert.Nil(t, repositories.Group.DropAll())
	assert.Nil(t, repositories.Organization.DropAll())
}

func TestCannotJoinChannelOutsideOrganization(t *testing.T) {
	_, user2, token := makeTestUsersAndToken(t)
	organization := makeTestOrganizations(t, []*models.User{user2})

	channel := makeTestChannel(organization.ID, false)

	w := httptest.NewRecorder()
	json := []byte(`{"channel_id":"` + channel.ID.Hex() + `"}`)
	req, _ := http.NewRequest("POST", "/api/channels/join", bytes.NewBuffer(json))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	assert.Nil(t, repositories.User.DropAll())
	assert.Nil(t, repositories.Group.DropAll())
	assert.Nil(t, repositories.Organization.DropAll())
}

func TestCannotJoinChannelTwice(t *testing.T) {
	user1, user2, token := makeTestUsersAndToken(t)
	organization := makeTestOrganizations(t, []*models.User{user1, user2})

	channel := makeTestChannel(organization.ID, false)

	w := httptest.NewRecorder()
	json := []byte(`{"channel_id":"` + channel.ID.Hex() + `"}`)
	req, _ := http.NewRequest("POST", "/api/channels/join", bytes.NewBuffer(json))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

	assert.Nil(t, repositories.User.DropAll())
	assert.Nil(t, repositories.Group.DropAll())
	assert.Nil(t, repositories.Organization.DropAll())
}

func TestCanLeaveChannel(t *testing.T) {
	user1, user2, token := makeTestUsersAndToken(t)
	organization := makeTestOrganizations(t, []*models.User{user1, user2})

	channel := makeTestChannelWithUser(organization.ID, false, user1)

	w := httptest.NewRecorder()
	json := []byte(`{"channel_id":"` + channel.ID.Hex() + `"}`)
	req, _ := http.NewRequest("POST", "/api/channels/leave", bytes.NewBuffer(json))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Nil(t, repositories.User.DropAll())
	assert.Nil(t, repositories.Group.DropAll())
	assert.Nil(t, repositories.Organization.DropAll())
}

func TestCannotLeaveInvalidChannel(t *testing.T) {
	_, _, token := makeTestUsersAndToken(t)

	w := httptest.NewRecorder()
	json := []byte(`{"channel_id":"non-existing"}`)
	req, _ := http.NewRequest("POST", "/api/channels/leave", bytes.NewBuffer(json))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	assert.Nil(t, repositories.User.DropAll())
}

func TestCannotLeaveChannelNotIn(t *testing.T) {
	user1, user2, token := makeTestUsersAndToken(t)
	organization := makeTestOrganizations(t, []*models.User{user1, user2})

	channel := makeTestChannelWithUser(organization.ID, false, user2)

	w := httptest.NewRecorder()
	json := []byte(`{"channel_id":"` + channel.ID.Hex() + `"}`)
	req, _ := http.NewRequest("POST", "/api/channels/leave", bytes.NewBuffer(json))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	assert.Nil(t, repositories.User.DropAll())
	assert.Nil(t, repositories.Group.DropAll())
	assert.Nil(t, repositories.Organization.DropAll())
}

func makeTestUsersAndToken(t *testing.T) (*models.User, *models.User, string) {
	user1, err := repositories.User.CreateUser("test1", "")
	assert.Nil(t, err)
	user2, err := repositories.User.CreateUser("test2", "")
	assert.Nil(t, err)
	token, _ := utils.CreateToken(user1)
	return user1, user2, token
}

func makeTestChannel(organizationID primitive.ObjectID, private bool) *models.Group {
	channelName := "test-channel"
	channelType := models.GroupTypePublicChannel
	if private {
		channelType = models.GroupTypePrivateMessage
	}
	group, _ := repositories.Group.CreateGroup(&models.Group{
		Name:         &channelName,
		Type:         channelType,
		Organization: organizationID,
	})
	return group
}

func makeTestChannelWithUser(organizationID primitive.ObjectID, private bool, user *models.User) *models.Group {
	group := makeTestChannel(organizationID, private)
	group.Users = append(group.Users, models.GroupUser{
		ID:       user.ID,
		Username: user.Username,
	})
	_ = repositories.Group.UpdateGroup(group)
	return group
}
