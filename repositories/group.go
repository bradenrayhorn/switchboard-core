package repositories

import (
	"github.com/Kamva/mgm/v3"
	"github.com/bradenrayhorn/switchboard-core/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GroupRepository interface {
	CreateGroup(group *models.Group) (*models.Group, error)
	GetGroups(userId primitive.ObjectID) ([]models.Group, error)
	ExistsByName(groupName string, organizationID primitive.ObjectID) (bool, error)
	GetByID(groupID primitive.ObjectID) (*models.Group, error)
	UpdateGroup(group *models.Group) error
	DropAll() error
}

var Group GroupRepository

func init() {
	Group = MongoGroupRepository{}
}

type MongoGroupRepository struct{}

func (m MongoGroupRepository) CreateGroup(group *models.Group) (*models.Group, error) {
	err := mgm.Coll(group).Create(group)
	return group, err
}

func (m MongoGroupRepository) GetGroups(userID primitive.ObjectID) ([]models.Group, error) {
	var groups = make([]models.Group, 0)
	cursor, err := mgm.Coll(&models.Group{}).Find(mgm.Ctx(), bson.M{"users.id": userID})
	if err != nil {
		return groups, err
	}

	err = cursor.All(mgm.Ctx(), &groups)

	return groups, nil
}

func (m MongoGroupRepository) ExistsByName(groupName string, organizationID primitive.ObjectID) (bool, error) {
	cursor, err := mgm.Coll(&models.Group{}).Find(mgm.Ctx(), bson.M{
		"name":         groupName,
		"organization": organizationID,
	})

	if err != nil {
		return false, err
	}
	cursor.Next(mgm.Ctx())

	return cursor.Current != nil, cursor.Close(mgm.Ctx())
}

func (m MongoGroupRepository) GetByID(groupID primitive.ObjectID) (*models.Group, error) {
	group := &models.Group{}
	err := mgm.Coll(&models.Group{}).First(bson.M{"_id": groupID}, group)

	if err != nil {
		return nil, err
	}

	return group, err
}

func (m MongoGroupRepository) UpdateGroup(group *models.Group) error {
	return mgm.Coll(group).Update(group)
}

func (m MongoGroupRepository) DropAll() error {
	return mgm.Coll(&models.Group{}).Drop(mgm.Ctx())
}
