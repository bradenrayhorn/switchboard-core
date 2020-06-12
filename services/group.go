package services

import (
	"github.com/bradenrayhorn/switchboard-backend/repositories"
	"github.com/bradenrayhorn/switchboard-backend/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
)

func CreateGroup(name string, userIds []string, authUserId string) *utils.HttpError {
	var primitiveUserIds []primitive.ObjectID
	var includesAuth = false
	for _, userId := range userIds {
		primitiveUserId, err := primitive.ObjectIDFromHex(userId)
		if err != nil {
			return utils.MakeHttpError(http.StatusUnprocessableEntity, "invalid user id provided")
		}
		if userId == authUserId {
			includesAuth = true
		}
		primitiveUserIds = append(primitiveUserIds, primitiveUserId)
	}

	if !includesAuth {
		return utils.MakeHttpError(http.StatusUnprocessableEntity, "you must be a member of the group")
	}

	exists, err := repositories.Group.GroupExists(primitiveUserIds)
	if err != nil {
		log.Println(err)
		return utils.MakeHttpError(http.StatusInternalServerError, "failed to create group")
	}

	if exists {
		return utils.MakeHttpError(http.StatusUnprocessableEntity, "group already exists")
	}

	var groupName *string = nil
	if len(name) > 0 {
		groupName = &name
	}

	_, err = repositories.Group.CreateGroup(groupName, primitiveUserIds)

	if err != nil {
		log.Println(err)
		return utils.MakeHttpError(http.StatusInternalServerError, "failed to create group")
	}

	return nil
}
