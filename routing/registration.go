package routing

import (
	"github.com/bradenrayhorn/switchboard-backend/repositories"
	"github.com/bradenrayhorn/switchboard-backend/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

func Register(c *gin.Context) {
	var request RegisterRequest
	if err := c.ShouldBind(&request); err != nil {
		utils.JsonError(http.StatusUnprocessableEntity, err.Error(), c)
		return
	}

	exists, err := repositories.User.Exists(request.Username)

	if err != nil || exists {
		utils.JsonError(http.StatusUnprocessableEntity, "username already exists", c)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		utils.JsonError(http.StatusUnprocessableEntity, "failed to create user", c)
		return
	}

	user, err := repositories.User.CreateUser(request.Username, string(hashedPassword))

	if err != nil {
		log.Println(err)
		utils.JsonError(http.StatusInternalServerError, "failed to create user", c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": user.ID,
	})
}
