package routing

import (
	"github.com/bradenrayhorn/switchboard-backend/middleware"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
)

func MakeRouter() *gin.Engine {
	router := gin.Default()
	log.Println("MAKE ROUTER")
	applyRoutes(router)
	return router
}

func MakeTestRouter() *gin.Engine {
	gin.DefaultWriter = ioutil.Discard
	gin.DefaultErrorWriter = ioutil.Discard
	router := gin.New()
	router.Use(gin.Recovery())
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

	api.POST("/groups/create", CreateGroup)
	api.GET("/groups", GetGroups)
}
