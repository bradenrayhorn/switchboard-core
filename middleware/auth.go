package middleware

import (
	"github.com/bradenrayhorn/switchboard-backend/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func getToken(header string) string {
	parts := strings.Split(header, " ")
	if len(parts) != 2 {
		return ""
	}
	return parts[1]
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := getToken(c.GetHeader("Authorization"))
		if len(tokenString) == 0 {
			tokenString, _ = c.GetQuery("auth")
			tokenString, _ = url.QueryUnescape(tokenString)
		}

		token, err := utils.ParseToken(tokenString)

		if err != nil {
			log.Println(err.Error())
			utils.JsonError(http.StatusUnauthorized, "invalid api token", c)
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("user_id", claims["user_id"])
			c.Set("user_username", claims["user_username"])
			c.Next()
		} else {
			utils.JsonError(http.StatusUnauthorized, "invalid api token", c)
			c.Abort()
			return
		}
	}
}
