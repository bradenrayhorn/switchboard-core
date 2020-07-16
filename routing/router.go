package routing

import (
	"github.com/bradenrayhorn/switchboard-core/database"
	"github.com/bradenrayhorn/switchboard-core/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

func MakeRouter(redis *database.RedisDB) *gin.Engine {
	router := gin.Default()
	router.Use(middleware.RedisMiddleware(redis))
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

	api.GET("/channels", GetChannels)
	api.POST("/channels", CreateChannel)
	api.POST("/channels/leave", LeaveChannel)
	api.POST("/channels/join", JoinChannel)

	api.GET("/organizations", GetOrganizations)
	api.POST("/organizations", CreateOrganization)
	api.POST("/organizations/invite-user", AddUserToOrganization)

	api.POST("/users/search", SearchUsers)
}
