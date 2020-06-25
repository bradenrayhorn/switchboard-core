package routing

import (
	"github.com/bradenrayhorn/switchboard-core/middleware"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

func MakeRouter() *gin.Engine {
	router := gin.Default()
	applyRoutes(router)
	return router
}

func MakeTestRouter() *gin.Engine {
	gin.DefaultWriter = ioutil.Discard
	gin.DefaultErrorWriter = ioutil.Discard
	router := gin.New()
	applyRoutes(router)
	return router
}

func applyRoutes(router *gin.Engine) {
	router.GET("/health-check", func(context *gin.Context) {
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
}
