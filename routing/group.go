package routing

import (
	"github.com/bradenrayhorn/switchboard-core/repositories"
	"github.com/bradenrayhorn/switchboard-core/services"
	"github.com/bradenrayhorn/switchboard-core/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
)

func CreateGroup(c *gin.Context) {
	var request CreateGroupRequest
	if err := c.ShouldBind(&request); err != nil {
		utils.JsonError(http.StatusUnprocessableEntity, err.Error(), c)
		return
	}

	err := services.CreateGroup(request.GroupName, request.UserIds, c.GetString("user_id"))

	if err != nil {
		utils.JsonError(err.Code, err.Error.Error(), c)
	}
}

func GetGroups(c *gin.Context) {
	userId, err := primitive.ObjectIDFromHex(c.GetString("user_id"))

	if err != nil {
		log.Println(err.Error())
		utils.JsonError(http.StatusInternalServerError, "failed to get groups", c)
		return
	}

	groups, err := repositories.Group.GetGroups(userId)

	if err != nil {
		log.Println(err.Error())
		utils.JsonError(http.StatusInternalServerError, "failed to get groups", c)
		return
	}

	c.JSON(http.StatusOK, utils.Json{"data": groups})
}
