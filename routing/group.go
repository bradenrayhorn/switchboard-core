package routing

import (
	"github.com/bradenrayhorn/switchboard-core/models"
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

func UpdateGroup(c *gin.Context) {
	var request UpdateGroupRequest
	if err := c.ShouldBind(&request); err != nil {
		utils.JsonError(http.StatusUnprocessableEntity, err.Error(), c)
		return
	}

	err := services.UpdateGroup(request.Id, c.GetString("user_id"), request.GroupName, request.AddUserIds, request.RemoveUserIds)

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

	// get users
	var userIDs []primitive.ObjectID
	for _, group := range groups {
		for _, id := range group.UserIds {
			userIDs = append(userIDs, id)
		}
	}

	usersList, err := repositories.User.GetUsers(userIDs)
	if err != nil {
		utils.JsonError(http.StatusInternalServerError, "failed to get users", c)
		return
	}

	var users = make(map[string]models.User)
	for _, user := range usersList {
		users[user.ID.Hex()] = user
	}

	// transform groups
	groupResponse := []models.GroupResponse{}
	for i := range groups {
		group := groups[i]
		var userResponse []models.GroupUser
		for _, userId := range group.UserIds {
			userResponse = append(userResponse, models.GroupUser{
				ID:   userId.Hex(),
				Name: users[userId.Hex()].Username,
			})
		}
		groupResponse = append(groupResponse, models.GroupResponse{
			ID:    group.ID,
			Name:  group.Name,
			Users: userResponse,
		})
	}

	c.JSON(http.StatusOK, utils.Json{"data": groupResponse})
}
