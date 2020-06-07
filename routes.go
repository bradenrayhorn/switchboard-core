package main

import (
	"github.com/bradenrayhorn/switchboard-backend/controllers"
	"github.com/bradenrayhorn/switchboard-backend/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

func registerRoutes(router *gin.Engine) {

	router.GET("/health-check", func(context *gin.Context) {
		context.String(http.StatusOK, "ok")
	})

	auth := router.Group("/api/auth")
	api := router.Group("/api")
	api.Use(middleware.AuthMiddleware())

	auth.POST("/login", controllers.Login)
	auth.POST("/register", controllers.Register)

	api.GET("/me", controllers.ShowMe)
}
