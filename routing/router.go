package routing

import (
	"github.com/bradenrayhorn/switchboard-core/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

func MakeRouter() *gin.Engine {
	router := gin.Default()
	applyRoutes(router)
	return router
}

func applyRoutes(router *gin.Engine) {
	router.GET("/api/health-check", func(context *gin.Context) {
		context.String(http.StatusOK, "ok")
	})

	auth := router.Group("/api/auth")
	api := router.Group("/api")
	api.Use(middleware.AuthMiddleware())

	auth.POST("/login", Login)
	auth.POST("/register", Register)

	api.GET("/me", ShowMe)

	api.GET("/groups", GetGroups)
	api.POST("/groups/create", CreateGroup)
	api.POST("/groups/update", UpdateGroup)

	api.POST("/users/search", SearchUsers)
}
