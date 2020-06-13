package routing

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
	r := MakeTestRouter()
	utils.SetupTestRsaKeys()

	repositories.Group = &repositories.MockGroupRepository{}
	repositories.User = &repositories.MockUserRepository{}

	user1, _ := repositories.User.CreateUser("george", "")
	user2, _ := repositories.User.CreateUser("test", "$2a$10$naqzJWUaOFm1/512Od.wPO4H8Vh8K38IGAb7rtgFizSflLVhpgMRG")

	token, _ := utils.CreateToken(user1)

	w := httptest.NewRecorder()
	form := url.Values{"users": []string{user1.ID.Hex(), user2.ID.Hex()}}
	req, _ := http.NewRequest("POST", "/api/groups/create", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	groupExists, err := repositories.Group.GroupExists([]primitive.ObjectID{user1.ID, user2.ID})

	assert.True(t, groupExists)
	assert.Nil(t, err)
}

func TestCantCreateDuplicateGroup(t *testing.T) {
	r := MakeTestRouter()
	utils.SetupTestRsaKeys()

	repositories.Group = &repositories.MockGroupRepository{}
	repositories.User = &repositories.MockUserRepository{}

	user1, _ := repositories.User.CreateUser("george", "")
	user2, _ := repositories.User.CreateUser("test", "$2a$10$naqzJWUaOFm1/512Od.wPO4H8Vh8K38IGAb7rtgFizSflLVhpgMRG")

	token, _ := utils.CreateToken(user1)

	_, _ = repositories.Group.CreateGroup(nil, []primitive.ObjectID{user1.ID, user2.ID})

	w := httptest.NewRecorder()
	form := url.Values{"users": []string{user2.ID.Hex(), user1.ID.Hex()}}
	req, _ := http.NewRequest("POST", "/api/groups/create", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestCantCreateGroupWithoutMe(t *testing.T) {
	r := MakeTestRouter()
	utils.SetupTestRsaKeys()

	repositories.Group = &repositories.MockGroupRepository{}
	repositories.User = &repositories.MockUserRepository{}

	user1, _ := repositories.User.CreateUser("george", "")
	user2, _ := repositories.User.CreateUser("test", "$2a$10$naqzJWUaOFm1/512Od.wPO4H8Vh8K38IGAb7rtgFizSflLVhpgMRG")

	token, _ := utils.CreateToken(user1)

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

	repositories.Group = &repositories.MockGroupRepository{}
	repositories.User = &repositories.MockUserRepository{}

	user1, _ := repositories.User.CreateUser("george", "")

	token, _ := utils.CreateToken(user1)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/groups/create", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

type GetGroupsResponse struct {
	Data []models.Group
}

func TestGetGroups(t *testing.T) {
	r := MakeTestRouter()
	utils.SetupTestRsaKeys()

	repositories.Group = &repositories.MockGroupRepository{}
	repositories.User = &repositories.MockUserRepository{}
	user1, _ := repositories.User.CreateUser("test", "$2a$10$naqzJWUaOFm1/512Od.wPO4H8Vh8K38IGAb7rtgFizSflLVhpgMRG")
	user2, _ := repositories.User.CreateUser("george", "")

	_, _ = repositories.Group.CreateGroup(nil, []primitive.ObjectID{user1.ID, user2.ID})
	_, _ = repositories.Group.CreateGroup(nil, []primitive.ObjectID{user1.ID})

	token, _ := utils.CreateToken(user1)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/groups", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var body GetGroupsResponse
	_ = json.Unmarshal(w.Body.Bytes(), &body)
	assert.Len(t, body.Data, 2)
}
