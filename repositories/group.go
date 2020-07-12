package repositories

import (
	"github.com/Kamva/mgm/v3"
	"github.com/Kamva/mgm/v3/operator"
	"github.com/bradenrayhorn/switchboard-core/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GroupRepository interface {
	CreateGroup(groupName *string, userIds []primitive.ObjectID, organizationID primitive.ObjectID) (*models.Group, error)
	GetGroups(userId primitive.ObjectID) ([]models.Group, error)
	GroupExists(userIds []primitive.ObjectID, organizationID primitive.ObjectID) (bool, error)
	GetGroup(groupId primitive.ObjectID, userId primitive.ObjectID) (*models.Group, error)
	UpdateGroup(group *models.Group) error
	DeleteGroup(group *models.Group) error
	DropAll() error
}

var Group GroupRepository

func init() {
	Group = MongoGroupRepository{}
}

type MongoGroupRepository struct{}

func (m MongoGroupRepository) CreateGroup(groupName *string, userIds []primitive.ObjectID, organizationID primitive.ObjectID) (*models.Group, error) {
	group := &models.Group{
		Name:         groupName,
		UserIds:      userIds,
		Organization: organizationID,
	}

	err := mgm.Coll(group).Create(group)
	return group, err
}

func (m MongoGroupRepository) GetGroups(userId primitive.ObjectID) ([]models.Group, error) {
	var groups = make([]models.Group, 0)
	cursor, err := mgm.Coll(&models.Group{}).Find(mgm.Ctx(), bson.M{"users": userId})
	if err != nil {
		return groups, err
	}

	err = cursor.All(mgm.Ctx(), &groups)

	return groups, nil
}

func (m MongoGroupRepository) GroupExists(userIds []primitive.ObjectID, organizationID primitive.ObjectID) (bool, error) {
	cursor, err := mgm.Coll(&models.Group{}).Find(mgm.Ctx(), bson.M{
		"users":        bson.M{operator.All: userIds, operator.Size: len(userIds)},
		"organization": organizationID,
	})

	if err != nil {
		return false, err
	}
	cursor.Next(mgm.Ctx())

	return cursor.Current != nil, cursor.Close(mgm.Ctx())
}

func (m MongoGroupRepository) GetGroup(groupId primitive.ObjectID, userId primitive.ObjectID) (*models.Group, error) {
	group := &models.Group{}
	err := mgm.Coll(&models.Group{}).First(bson.M{"users": userId, "_id": groupId}, group)

	if err != nil {
		return nil, err
	}

	return group, err
}

func (m MongoGroupRepository) UpdateGroup(group *models.Group) error {
	return mgm.Coll(group).Update(group)
}

func (m MongoGroupRepository) DeleteGroup(group *models.Group) error {
	return mgm.Coll(group).Delete(group)
}

func (m MongoGroupRepository) DropAll() error {
	return mgm.Coll(&models.Group{}).Drop(mgm.Ctx())
}
