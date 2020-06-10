package routing

import (
	"github.com/bradenrayhorn/switchboard-backend/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

func MakeRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/health-check", func(context *gin.Context) {
		context.String(http.StatusOK, "ok")
	})

	auth := router.Group("/api/auth")
	api := router.Group("/api")
	api.Use(middleware.AuthMiddleware())

	auth.POST("/login", Login)
	auth.POST("/register", Register)

	api.GET("/me", ShowMe)

	return router
}
