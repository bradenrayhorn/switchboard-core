package services

import (
	"github.com/bradenrayhorn/switchboard-core/models"
	"github.com/bradenrayhorn/switchboard-core/repositories"
	"github.com/bradenrayhorn/switchboard-core/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"regexp"
	"strings"
)

func CreateChannel(channelName string, isPrivate bool, organizationID primitive.ObjectID, userID primitive.ObjectID, username string) (*models.Group, *utils.HttpError) {
	// validate organization
	organization, err := repositories.Organization.GetForUserAndID(organizationID, userID)
	if err != nil || organization == nil {
		return nil, utils.MakeHttpError(http.StatusUnprocessableEntity, "you do not have access to that organization")
	}

	// clean and validate name
	cleanedName := strings.TrimSpace(channelName)
	if !regexp.MustCompile(`^[a-zA-Z_-]+$`).MatchString(cleanedName) {
		return nil, utils.MakeHttpError(http.StatusUnprocessableEntity, "name is invalid")
	}

	// validate name is unique
	exists, err := repositories.Group.ExistsByName(cleanedName, organizationID)
	if err != nil || exists {
		return nil, utils.MakeHttpError(http.StatusUnprocessableEntity, "channel already exists")
	}

	// create channel
	channelType := models.GroupTypePublicChannel
	if isPrivate {
		channelType = models.GroupTypePrivateChannel
	}

	group := models.Group{
		Name: &channelName,
		Type: channelType,
		Users: []models.GroupUser{{
			ID:       userID,
			Username: username,
		}},
		Organization: organizationID,
	}

	_, err = repositories.Group.CreateGroup(&group)
	if err != nil {
		return nil, utils.MakeHttpError(http.StatusInternalServerError, "failed to create channel")
	}

	return &group, nil
}

func JoinChannel(channelID primitive.ObjectID, userID primitive.ObjectID, username string) *utils.HttpError {
	// find group
	group, err := repositories.Group.GetByID(channelID)
	if err != nil || group == nil || group.Type == models.GroupTypePrivateChannel || group.Type == models.GroupTypePrivateMessage {
		return utils.MakeHttpError(http.StatusUnprocessableEntity, "could not find channel")
	}

	// verify group organization
	organization, err := repositories.Organization.GetForUserAndID(group.Organization, userID)
	if err != nil || organization == nil {
		return utils.MakeHttpError(http.StatusUnprocessableEntity, "you do not have access to that channel")
	}

	// verify user is not in group
	for _, user := range group.Users {
		if user.ID == userID {
			return utils.MakeHttpError(http.StatusUnprocessableEntity, "you already joined the channel")
		}
	}

	// add user to group
	group.Users = append(group.Users, models.GroupUser{
		ID:       userID,
		Username: username,
	})

	err = repositories.Group.UpdateGroup(group)
	if err != nil {
		return utils.MakeHttpError(http.StatusInternalServerError, "failed to join channel")
	}
	return nil
}

func LeaveChannel(channelID primitive.ObjectID, userID primitive.ObjectID) *utils.HttpError {
	// find group
	group, err := repositories.Group.GetByID(channelID)
	if err != nil || group == nil {
		return utils.MakeHttpError(http.StatusUnprocessableEntity, "could not find channel")
	}

	// verify user is in group
	userFound := false
	for _, user := range group.Users {
		if user.ID == userID {
			userFound = true
		}
	}
	if !userFound {
		return utils.MakeHttpError(http.StatusUnprocessableEntity, "you are not in the channel")
	}

	// remove user from group
	for i, user := range group.Users {
		if user.ID == userID {
			group.Users = append(group.Users[:i], group.Users[i+1:]...)
			break
		}
	}

	err = repositories.Group.UpdateGroup(group)
	if err != nil {
		return utils.MakeHttpError(http.StatusInternalServerError, "failed to leave channel")
	}
	return nil
}

func GetChannelsForOrganization(organizationID primitive.ObjectID, userID primitive.ObjectID) ([]models.Group, *utils.HttpError) {
	// verify user is in organization
	organization, err := repositories.Organization.GetForUserAndID(organizationID, userID)
	if err != nil || organization == nil {
		return nil, utils.MakeHttpError(http.StatusUnprocessableEntity, "you do not have access to that organization")
	}

	// get channels for organization
	groups, err := repositories.Group.GetByOrganization(organizationID, models.GroupTypePublicChannel)
	if err != nil {
		return nil, utils.MakeHttpError(http.StatusInternalServerError, "failed to get channels")
	}
	return groups, nil
}
