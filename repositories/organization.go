package repositories

import (
	"github.com/Kamva/mgm/v3"
	"github.com/bradenrayhorn/switchboard-core/models"
)

type OrganizationRepository interface {
	Create(name string, users []models.OrganizationUser) (*models.Organization, error)
	DropAll() error
}

var Organization OrganizationRepository

func init() {
	Organization = MongoOrganizationRepository{}
}

type MongoOrganizationRepository struct{}

func (m MongoOrganizationRepository) Create(name string, users []models.OrganizationUser) (*models.Organization, error) {
	organization := &models.Organization{
		Name:  name,
		Users: users,
	}

	err := mgm.Coll(organization).Create(organization)
	return organization, err
}

func (m MongoOrganizationRepository) DropAll() error {
	return mgm.Coll(&models.Organization{}).Drop(mgm.Ctx())
}
