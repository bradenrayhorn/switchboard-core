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

func GetOrganizations(userID string) ([]models.Organization, *utils.HttpError) {
	var organizations = make([]models.Organization, 0)
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return organizations, utils.MakeHttpError(http.StatusUnprocessableEntity, "invalid user id")
	}

	organizations, err = repositories.Organization.GetForUser(userObjectID)
	if err != nil {
		log.Println(err)
		return nil, utils.MakeHttpError(http.StatusInternalServerError, "failed to get organizations")
	}

	return organizations, nil
}

func AddUserToOrganization(organizationID string, username string, authUserID string) *utils.HttpError {
	// prepare ids
	userObjectID, userErr := primitive.ObjectIDFromHex(authUserID)
	organizationObjectID, orgErr := primitive.ObjectIDFromHex(organizationID)
	if userErr != nil || orgErr != nil {
		return utils.MakeHttpError(http.StatusUnprocessableEntity, "invalid id")
	}

	// check permission
	organization, err := repositories.Organization.GetForUserAndID(organizationObjectID, userObjectID)
	if err != nil || organization == nil {
		return utils.MakeHttpError(http.StatusUnprocessableEntity, "invalid permission to invite to organization")
	}

	// find user
	user, err := repositories.User.GetUser(username)
	if err != nil {
		return utils.MakeHttpError(http.StatusUnprocessableEntity, "could not find user")
	}

	// ensure user not already in organization
	for _, orgUser := range organization.Users {
		if orgUser.ID == user.ID {
			return utils.MakeHttpError(http.StatusUnprocessableEntity, "user already in organization")
		}
	}

	// add user to organization
	userRole := models.OrganizationUser{
		ID:   user.ID,
		Role: models.RoleUser,
	}
	organization.Users = append(organization.Users, userRole)
	err = repositories.Organization.UpdateOrganization(organization)
	if err != nil {
		return utils.MakeHttpError(http.StatusInternalServerError, "failed to save organization")
	} else {
		return nil
	}
}
