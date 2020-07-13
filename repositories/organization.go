package repositories

import (
	"github.com/Kamva/mgm/v3"
	"github.com/Kamva/mgm/v3/operator"
	"github.com/bradenrayhorn/switchboard-core/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrganizationRepository interface {
	Create(name string, users []models.OrganizationUser) (*models.Organization, error)
	GetForUser(userID primitive.ObjectID) ([]models.Organization, error)
	GetForUserAndID(organizationID primitive.ObjectID, userID primitive.ObjectID) (*models.Organization, error)
	UpdateOrganization(organization *models.Organization) error
	FindUser(organizationID primitive.ObjectID, username string) ([]models.OrganizationUser, error)
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

func (m MongoOrganizationRepository) GetForUser(userID primitive.ObjectID) ([]models.Organization, error) {
	var organizations = make([]models.Organization, 0)
	cursor, err := mgm.Coll(&models.Organization{}).Find(mgm.Ctx(), bson.M{"users.id": userID})
	if err != nil {
		return organizations, err
	}

	err = cursor.All(mgm.Ctx(), &organizations)

	return organizations, nil
}

func (m MongoOrganizationRepository) GetForUserAndID(organizationID primitive.ObjectID, userID primitive.ObjectID) (*models.Organization, error) {
	organization := &models.Organization{}
	err := mgm.Coll(&models.Organization{}).First(bson.M{"users.id": userID, "_id": organizationID}, organization)

	if err != nil {
		return nil, err
	}

	return organization, err
}

func (m MongoOrganizationRepository) UpdateOrganization(organization *models.Organization) error {
	return mgm.Coll(organization).Update(organization)
}

func (m MongoOrganizationRepository) FindUser(organizationID primitive.ObjectID, username string) ([]models.OrganizationUser, error) {
	users := make([]models.OrganizationUser, 0)
	err := mgm.Coll(&models.Organization{}).SimpleAggregate(&users,
		bson.M{operator.Match: bson.M{"_id": organizationID}},
		bson.M{operator.Unwind: "$users"},
		bson.M{operator.Match: bson.M{"users.username": bson.M{operator.Regex: username}}},
		bson.M{operator.Project: bson.M{"id": "$users.id", "username": "$users.username"}},
	)
	return users, err
}

func (m MongoOrganizationRepository) DropAll() error {
	return mgm.Coll(&models.Organization{}).Drop(mgm.Ctx())
}
