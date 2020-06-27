package tests

import (
	"github.com/Kamva/mgm/v3"
	"github.com/bradenrayhorn/switchboard-core/models"
	"github.com/bradenrayhorn/switchboard-core/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func TestCreateToken(t *testing.T) {
	// create test user
	user := models.User{
		DefaultModel: mgm.DefaultModel{
			IDField:    mgm.IDField{ID: primitive.NewObjectID()},
			DateFields: mgm.DateFields{},
		},
		Username: "test",
		Password: "",
	}

	// create token
	token, err := utils.CreateToken(&user)

	require.Nil(t, err, "creating token should succeed")
	assert.True(t, len(token) > 0, "token has content")
}

func TestParseToken(t *testing.T) {
	// create test user
	userId := primitive.NewObjectID()
	user := models.User{
		DefaultModel: mgm.DefaultModel{
			IDField:    mgm.IDField{ID: userId},
			DateFields: mgm.DateFields{},
		},
		Username: "test",
		Password: "",
	}

	// create token
	tokenString, err := utils.CreateToken(&user)

	require.Nil(t, err, "creating token should succeed")
	assert.True(t, len(tokenString) > 0, "token has content")

	// parse token
	token, err := utils.ParseToken(tokenString)

	require.Nil(t, err, "parsing token should succeed")

	claims, ok := token.Claims.(jwt.MapClaims)

	assert.True(t, ok, "claims can be parsed")
	assert.True(t, token.Valid, "token is valid")
	assert.Equal(t, "test", claims["user_username"], "username matches")
	assert.Equal(t, userId.Hex(), claims["user_id"], "user id matches")
}

func TestParseTokenFail(t *testing.T) {
	token := jwt.New(jwt.SigningMethodHS256)

	tokenString, err := token.SignedString([]byte("my-jwt-secret"))

	require.Nil(t, err)

	_, err = utils.ParseToken(tokenString)

	assert.NotNil(t, err)
}
