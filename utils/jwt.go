package utils

import (
	"fmt"
	"github.com/bradenrayhorn/switchboard-core/config"
	"github.com/bradenrayhorn/switchboard-core/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"time"
)

func CreateToken(user *models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"user_username": user.Username,
		"user_id":       user.ID.Hex(),
		"exp":           time.Now().Add(viper.GetDuration("token_expiration")).Unix(),
	})
	return token.SignedString(config.RsaPrivate)
}

func ParseToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return config.RsaPublic, nil
	})
}
