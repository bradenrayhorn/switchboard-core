package routing

import (
	"github.com/bradenrayhorn/switchboard-core/services"
	"github.com/bradenrayhorn/switchboard-core/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateOrganization(c *gin.Context) {
	var request CreateOrganizationRequest
	if err := c.ShouldBind(&request); err != nil {
		utils.JsonError(http.StatusUnprocessableEntity, err.Error(), c)
		return
	}

	_, err := services.CreateOrganization(request.Name, c.GetString("user_id"), c.GetString("user_username"))

	if err != nil {
		utils.JsonError(err.Code, err.Error.Error(), c)
		return
	}
}

func GetOrganizations(c *gin.Context) {
	organizations, err := services.GetOrganizations(c.GetString("user_id"))

	if err != nil {
		utils.JsonError(err.Code, err.Error.Error(), c)
	} else {
		c.JSON(http.StatusOK, utils.Json{"data": organizations})
	}
}

func AddUserToOrganization(c *gin.Context) {
	var request AddUserToOrganizationRequest
	if err := c.ShouldBind(&request); err != nil {
		utils.JsonError(http.StatusUnprocessableEntity, err.Error(), c)
		return
	}

	err := services.AddUserToOrganization(request.OrganizationID, request.Username, c.GetString("user_id"))
	if err != nil {
		utils.JsonError(err.Code, err.Error.Error(), c)
		return
	}
}
