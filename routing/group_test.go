package routing

import (
	"encoding/json"
	"fmt"
	"github.com/bradenrayhorn/switchboard-core/models"
	"github.com/bradenrayhorn/switchboard-core/repositories"
	"github.com/bradenrayhorn/switchboard-core/repositories/mocks"
	"github.com/bradenrayhorn/switchboard-core/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestCreateGroup(t *testing.T) {
	r := MakeTestRouter()
	utils.SetupTestRsaKeys()

	user1 := utils.MakeTestUser("test1", "")
	user2 := utils.MakeTestUser("test2", "")
	userIDSlice := []primitive.ObjectID{user1.ID, user2.ID}
	token, _ := utils.CreateToken(user1)

	groupRepo := new(mocks.GroupRepository)
	groupRepo.On("GroupExists", userIDSlice).Return(false, nil)
	groupRepo.On("CreateGroup", mock.AnythingOfType("*string"), userIDSlice).Return(nil, nil)
	repositories.Group = groupRepo

	w := httptest.NewRecorder()
	form := url.Values{"users": []string{user1.ID.Hex(), user2.ID.Hex()}}
	req, _ := http.NewRequest("POST", "/api/groups/create", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	groupRepo.AssertCalled(t, "CreateGroup", mock.AnythingOfType("*string"), userIDSlice)
}

func TestCantCreateDuplicateGroup(t *testing.T) {
	r := MakeTestRouter()
	utils.SetupTestRsaKeys()

	user1 := utils.MakeTestUser("test1", "")
	user2 := utils.MakeTestUser("test2", "")
	userIDSlice := []primitive.ObjectID{user2.ID, user1.ID}
	token, _ := utils.CreateToken(user1)

	groupRepo := new(mocks.GroupRepository)
	groupRepo.On("GroupExists", userIDSlice).Return(true, nil)
	repositories.Group = groupRepo

	w := httptest.NewRecorder()
	form := url.Values{"users": []string{user2.ID.Hex(), user1.ID.Hex()}}
	req, _ := http.NewRequest("POST", "/api/groups/create", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	groupRepo.AssertNotCalled(t, "CreateGroup", mock.AnythingOfType("*string"), userIDSlice)
}

func TestCantCreateGroupWithoutMe(t *testing.T) {
	r := MakeTestRouter()
	utils.SetupTestRsaKeys()

	user1 := utils.MakeTestUser("test1", "")
	user2 := utils.MakeTestUser("test2", "")
	token, _ := utils.CreateToken(user1)

	groupRepo := new(mocks.GroupRepository)
	repositories.Group = groupRepo

	w := httptest.NewRecorder()
	form := url.Values{"users": []string{user2.ID.Hex()}}
	req, _ := http.NewRequest("POST", "/api/groups/create", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestCantCreateGroupWithoutUsers(t *testing.T) {
	r := MakeTestRouter()
	utils.SetupTestRsaKeys()

	user1 := utils.MakeTestUser("test1", "")
	token, _ := utils.CreateToken(user1)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/groups/create", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestGetGroups(t *testing.T) {
	r := MakeTestRouter()
	utils.SetupTestRsaKeys()

	user1 := utils.MakeTestUser("test1", "")
	user2 := utils.MakeTestUser("test2", "")
	token, _ := utils.CreateToken(user1)

	// setup group repo
	groupRepo := new(mocks.GroupRepository)
	groupRepo.On("GetGroups", user1.ID).Return([]models.Group{
		utils.MakeTestGroup(nil, []primitive.ObjectID{user1.ID}),
		utils.MakeTestGroup(nil, []primitive.ObjectID{user1.ID, user2.ID}),
	}, nil)
	repositories.Group = groupRepo
	// setup user repo
	userRepo := new(mocks.UserRepository)
	userRepo.On("GetUsers", mock.Anything).Return([]models.User{*user1, *user2}, nil)
	repositories.User = userRepo

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
