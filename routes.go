package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func registerRoutes(router *gin.Engine) {

	router.GET("/health-check", func(context *gin.Context) {
		context.String(http.StatusOK, "ok")
	})

}
