package utils

import (
	"github.com/Kamva/mgm/v3"
	"github.com/bradenrayhorn/switchboard-backend/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func TestCreateToken(t *testing.T) {
	// set viper jwt secret
	viper.Set("jwt_secret", "my-jwt-secret")

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
	token, err := CreateToken(&user)

	require.Nil(t, err, "creating token should succeed")
	assert.True(t, len(token) > 0, "token has content")
}

func TestParseToken(t *testing.T) {
	// set viper jwt secret
	viper.Set("jwt_secret", "my-jwt-secret")

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
	tokenString, err := CreateToken(&user)

	require.Nil(t, err, "creating token should succeed")
	assert.True(t, len(tokenString) > 0, "token has content")

	// parse token
	token, err := ParseToken(tokenString)

	require.Nil(t, err, "parsing token should succeed")

	claims, ok := token.Claims.(jwt.MapClaims)

	assert.True(t, ok, "claims can be parsed")
	assert.True(t, token.Valid, "token is valid")
	assert.Equal(t, "test", claims["user_username"], "username matches")
	assert.Equal(t, userId.Hex(), claims["user_id"], "user id matches")
}
