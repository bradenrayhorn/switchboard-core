package services

import (
	"github.com/bradenrayhorn/switchboard-core/models"
	"github.com/bradenrayhorn/switchboard-core/repositories"
	"github.com/bradenrayhorn/switchboard-core/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
)

func CreateOrganization(name string, userID string) (*models.Organization, *utils.HttpError) {
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, utils.MakeHttpError(http.StatusUnprocessableEntity, "invalid user id")
	}

	users := []models.OrganizationUser{{
		ID:   userObjectID,
		Role: models.RoleAdmin,
	}}

	organization, err := repositories.Organization.Create(name, users)

	if err != nil {
		log.Println(err)
		return nil, utils.MakeHttpError(http.StatusInternalServerError, "failed to create organization")
	}

	return organization, nil
}
