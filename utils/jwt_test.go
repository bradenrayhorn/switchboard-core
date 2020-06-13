package utils

import (
	"github.com/Kamva/mgm/v3"
	"github.com/bradenrayhorn/switchboard-core/models"
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

const testRsaKey = `
-----BEGIN RSA PRIVATE KEY-----
MIICXAIBAAKBgQCXrkec36nQ3PdYcPc2v648fPJgUimbkls6Yy8XYHyTUaSdn8jB
POBhkJiY/ClH1JJCYwGTJQ2Zw6ud+m1HsVVLcagYAoUDLXeOvRYCfSWnLIoMAnt4
gpD+6ZIE2F5hAAmNWNpZ/vGGDcEbAA82JzuPsvrZwUhISP3lcBjvqFFR/QIDAQAB
AoGAFu8F6uUyct8GEvw5lLCUspaduwyRN/GAE6rtctZm34tnnWGMZBNFRdssB22Q
/EhauOjpUws2LDqWlzNHFKDuaVWr0N1NdI5l28fyq0r66oyHu1QaBhwGaGYqlQta
krVtJiGPjNVYgjAT7CbiPlBHQRYc8lKJ7I9Q35Ueijs8boECQQD3P0Hb45H8Elcy
uo97aY4sw9VHj9JI6SXRYKsn0GjoJkUeTO+BLZVRfyU6XVpjDehToEGKdOhCjcj+
w8gJKOHhAkEAnQztX6Pzqfuf0OAbM8w7HaHkrqPFRxo1hIUTkyoC8a5vWb5dAg7f
zFmeev7KUGwjtRLqsW7FSuPHCDTaEg4rnQJAdEmkSC+4sb5OlOV6Jql23UceQRI7
7V77EodM+OTk8otNZvr4iuFNoY8Ti2fs4znfU7XEHcLump7lMi4TH3vDAQJAECAI
q15BIX3lfePUxy+8DiHWQhwsoE1Pm2iyhtS0cq4vXq6ODC0My4WUIRrSGQnRbMEh
edBez36tq+kJTvtHfQJBANd4Dgp65qFDhPCaNZcYRW8KhjGk4NFHLBbHCuVrGXFG
nxI1y3andwCyT2w0ol/ci5Lw3SsZlWXcQRYiRgzqukQ=
-----END RSA PRIVATE KEY-----
`

func TestParseTokenFail(t *testing.T) {
	viper.Set("jwt_secret", "my-jwt-secret")

	token := jwt.New(jwt.GetSigningMethod("RS256"))
	key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(testRsaKey))

	require.Nil(t, err)

	tokenString, err := token.SignedString(key)

	require.Nil(t, err)

	_, err = ParseToken(tokenString)

	assert.NotNil(t, err)
}
