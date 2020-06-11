package utils

import (
	"fmt"
	"github.com/bradenrayhorn/switchboard-backend/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"time"
)

func CreateToken(user *models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_username": user.Username,
		"user_id":       user.ID.Hex(),
		"exp":           time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString([]byte(viper.GetString("jwt_secret")))
}

func ParseToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(viper.GetString("jwt_secret")), nil
	})
}
