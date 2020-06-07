package controllers

import (
	"github.com/bradenrayhorn/switchboard-backend/repositories"
	"github.com/bradenrayhorn/switchboard-backend/utils"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func Register(c *gin.Context) {
	var request RegisterRequest
	if err := c.ShouldBind(&request); err != nil {
		utils.JsonError(http.StatusUnprocessableEntity, err.Error(), c)
		return
	}

	exists, err := repositories.UserExists(request.Username)

	if err != nil || exists {
		utils.JsonError(http.StatusUnprocessableEntity, "username already exists", c)
		return
	}

	user, err := repositories.CreateUser(request.Username, request.Password)

	if err != nil {
		log.Println(err)
		utils.JsonError(http.StatusUnprocessableEntity, "failed to create user", c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": user.ID,
	})
}
