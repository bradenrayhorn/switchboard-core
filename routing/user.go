package routing

import (
	"github.com/bradenrayhorn/switchboard-core/models"
	"github.com/bradenrayhorn/switchboard-core/repositories"
	"github.com/bradenrayhorn/switchboard-core/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func SearchUsers(c *gin.Context) {
	var request UserSearchRequest
	if err := c.ShouldBind(&request); err != nil {
		utils.JsonError(http.StatusUnprocessableEntity, err.Error(), c)
		return
	}

	organizationID, err := primitive.ObjectIDFromHex(request.OrganizationID)
	if err != nil {
		utils.JsonError(http.StatusUnprocessableEntity, "invalid organization id", c)
		return
	}

	users, err := repositories.Organization.FindUser(organizationID, request.Name)

	if err != nil {
		utils.JsonError(http.StatusInternalServerError, "failed to search", c)
		return
	}

	userResponse := make([]models.UserResponse, 0)
	for i := range users {
		userResponse = append(userResponse, models.UserResponse{
			ID:       users[i].ID,
			Username: users[i].Username,
		})
	}

	c.JSON(http.StatusOK, utils.Json{"data": userResponse})
}
