package routing

import (
	"github.com/bradenrayhorn/switchboard-core/models"
	"github.com/bradenrayhorn/switchboard-core/repositories"
	"github.com/bradenrayhorn/switchboard-core/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SearchUsers(c *gin.Context) {
	var request UserSearchRequest
	if err := c.ShouldBind(&request); err != nil {
		utils.JsonError(http.StatusUnprocessableEntity, err.Error(), c)
		return
	}

	users, err := repositories.User.SearchUser(request.Name)

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
