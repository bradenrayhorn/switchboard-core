package services

import (
	"github.com/bradenrayhorn/switchboard-core/models"
	"github.com/bradenrayhorn/switchboard-core/repositories"
	"github.com/bradenrayhorn/switchboard-core/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
)

func CreateGroup(name string, userIds []string, authUserId string) (*models.Group, *utils.HttpError) {
	var primitiveUserIds []primitive.ObjectID
	var includesAuth = false
	for _, userId := range userIds {
		primitiveUserId, err := primitive.ObjectIDFromHex(userId)
		if err != nil {
			return nil, utils.MakeHttpError(http.StatusUnprocessableEntity, "invalid user id provided")
		}
		if userId == authUserId {
			includesAuth = true
		}
		primitiveUserIds = append(primitiveUserIds, primitiveUserId)
	}

	if !includesAuth {
		return nil, utils.MakeHttpError(http.StatusUnprocessableEntity, "you must be a member of the group")
	}

	exists, err := repositories.Group.GroupExists(primitiveUserIds)
	if err != nil {
		log.Println(err)
		return nil, utils.MakeHttpError(http.StatusInternalServerError, "failed to create group")
	}

	if exists {
		return nil, utils.MakeHttpError(http.StatusUnprocessableEntity, "group already exists")
	}

	var groupName *string = nil
	if len(name) > 0 {
		groupName = &name
	}

	group, err := repositories.Group.CreateGroup(groupName, primitiveUserIds)

	if err != nil {
		log.Println(err)
		return nil, utils.MakeHttpError(http.StatusInternalServerError, "failed to create group")
	}

	return group, nil
}

func UpdateGroup(id string, authUserId string, groupName string, usersToAdd []string, usersToRemove []string) *utils.HttpError {
	// parse ids
	groupId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return utils.MakeHttpError(http.StatusUnprocessableEntity, "invalid group id")
	}
	userId, err := primitive.ObjectIDFromHex(authUserId)
	if err != nil {
		return utils.MakeHttpError(http.StatusUnprocessableEntity, "invalid user id")
	}

	var userIdsToAdd []primitive.ObjectID
	for _, userId := range usersToAdd {
		primitiveUserId, err := primitive.ObjectIDFromHex(userId)
		if err != nil {
			return utils.MakeHttpError(http.StatusUnprocessableEntity, "invalid user id provided")
		}
		userIdsToAdd = append(userIdsToAdd, primitiveUserId)
	}

	var userIdsToRemove []primitive.ObjectID
	for _, userId := range usersToRemove {
		primitiveUserId, err := primitive.ObjectIDFromHex(userId)
		if err != nil {
			return utils.MakeHttpError(http.StatusUnprocessableEntity, "invalid user id provided")
		}
		userIdsToRemove = append(userIdsToRemove, primitiveUserId)
	}

	// get group
	group, err := repositories.Group.GetGroup(groupId, userId)
	if err != nil {
		return utils.MakeHttpError(http.StatusInternalServerError, "failed to update group")
	}

	if group == nil {
		return utils.MakeHttpError(http.StatusUnprocessableEntity, "you do not have access to that group")
	}

	// update name
	if len(groupName) > 0 {
		group.Name = &groupName
	} else {
		group.Name = nil
	}

	// add users
addLoop:
	for _, toAdd := range userIdsToAdd {
		for _, userId := range group.UserIds {
			if userId == toAdd {
				break addLoop
			}
		}
		group.UserIds = append(group.UserIds, toAdd)
	}

	// remove users
	for _, toRemove := range userIdsToRemove {
		index := -1
		for i, userId := range group.UserIds {
			if userId == toRemove {
				index = i
				break
			}
		}
		group.UserIds = append(group.UserIds[:index], group.UserIds[index+1:]...)
	}

	if len(group.UserIds) == 0 {
		err = repositories.Group.DeleteGroup(group)
	} else {
		err = repositories.Group.UpdateGroup(group)
	}

	if err != nil {
		return utils.MakeHttpError(http.StatusInternalServerError, "failed to update group")
	}

	return nil
}
