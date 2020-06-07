package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func ShowMe(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"id": c.GetString("user_id"),
	})
}
