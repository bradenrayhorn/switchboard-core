package routing

import (
	"github.com/bradenrayhorn/switchboard-core/database"
	"github.com/bradenrayhorn/switchboard-core/repositories"
	"github.com/bradenrayhorn/switchboard-core/services"
	"github.com/bradenrayhorn/switchboard-core/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func CreateChannel(c *gin.Context) {
	var request CreateChannelRequest
	if err := c.ShouldBind(&request); err != nil {
		utils.JsonError(http.StatusUnprocessableEntity, err.Error(), c)
		return
	}

	userID := c.MustGet("user_id_object").(primitive.ObjectID)
	group, err := services.CreateChannel(request.Name, request.Private, request.OrganizationID, userID, c.GetString("user_username"))

	if err != nil {
		utils.JsonError(err.Code, err.Error.Error(), c)
		return
	}

	redis := c.MustGet("redis").(*database.RedisDB)
	redis.PublishGroupJoin(userID.Hex(), group.ID.Hex())
}

func GetChannels(c *gin.Context) {
	userID := c.MustGet("user_id_object").(primitive.ObjectID)
	channels, err := repositories.Group.GetGroups(userID)
	if err != nil {
		utils.JsonError(http.StatusInternalServerError, "failed to get channels", c)
		return
	} else {
		c.JSON(http.StatusOK, utils.Json{"data": channels})
	}
}

func LeaveChannel(c *gin.Context) {
	var request LeaveChannelRequest
	if err := c.ShouldBind(&request); err != nil {
		utils.JsonError(http.StatusUnprocessableEntity, err.Error(), c)
		return
	}

	userID := c.MustGet("user_id_object").(primitive.ObjectID)
	err := services.LeaveChannel(request.ChannelID, userID)

	if err != nil {
		utils.JsonError(err.Code, err.Error.Error(), c)
		return
	}

	redis := c.MustGet("redis").(*database.RedisDB)
	redis.PublishGroupLeft(userID.Hex(), request.ChannelID.Hex())
}

func JoinChannel(c *gin.Context) {
	var request JoinChannelRequest
	if err := c.ShouldBind(&request); err != nil {
		utils.JsonError(http.StatusUnprocessableEntity, err.Error(), c)
		return
	}

	userID := c.MustGet("user_id_object").(primitive.ObjectID)
	err := services.JoinChannel(request.ChannelID, userID, c.GetString("user_username"))

	if err != nil {
		utils.JsonError(err.Code, err.Error.Error(), c)
		return
	}

	redis := c.MustGet("redis").(*database.RedisDB)
	redis.PublishGroupJoin(userID.Hex(), request.ChannelID.Hex())
}
