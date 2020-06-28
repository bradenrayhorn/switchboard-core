package routing

import (
	"github.com/bradenrayhorn/switchboard-core/repositories"
	"github.com/bradenrayhorn/switchboard-core/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

func Login(c *gin.Context) {
	var request LoginRequest
	if err := c.ShouldBind(&request); err != nil {
		utils.JsonError(http.StatusUnprocessableEntity, err.Error(), c)
		return
	}

	user, err := repositories.User.GetUser(request.Username)

	if err != nil {
		utils.JsonError(http.StatusUnprocessableEntity, "invalid username/password", c)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		utils.JsonError(http.StatusUnprocessableEntity, "invalid username/password", c)
		return
	}

	token, err := utils.CreateToken(user)
	if err != nil {
		log.Println(err.Error())
		utils.JsonError(http.StatusInternalServerError, "internal error", c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"id":    user.ID,
	})
}
