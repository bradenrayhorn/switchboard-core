package middleware

import (
	"github.com/bradenrayhorn/switchboard-core/database"
	"github.com/gin-gonic/gin"
)

func RedisMiddleware(redis *database.RedisDB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("redis", redis)
		c.Next()
	}
}
