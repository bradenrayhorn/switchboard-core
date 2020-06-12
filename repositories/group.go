package repositories

import (
	"github.com/Kamva/mgm/v3"
	"github.com/Kamva/mgm/v3/operator"
	"github.com/bradenrayhorn/switchboard-backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GroupRepository interface {
	CreateGroup(groupName *string, userIds []primitive.ObjectID) (*models.Group, error)
	GetGroups(userId primitive.ObjectID) ([]models.Group, error)
	GroupExists(userIds []primitive.ObjectID) (bool, error)
}

var Group GroupRepository

func init() {
	Group = MongoGroupRepository{}
}

type MongoGroupRepository struct{}

func (m MongoGroupRepository) CreateGroup(groupName *string, userIds []primitive.ObjectID) (*models.Group, error) {
	group := &models.Group{
		Name:    groupName,
		UserIds: userIds,
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

func (m MongoGroupRepository) GroupExists(userIds []primitive.ObjectID) (bool, error) {
	cursor, err := mgm.Coll(&models.Group{}).Find(mgm.Ctx(), bson.M{"users": bson.M{operator.All: userIds, operator.Size: len(userIds)}})

	if err != nil {
		return false, err
	}
	cursor.Next(mgm.Ctx())

	return cursor.Current != nil, cursor.Close(mgm.Ctx())
}

// mock repository

type MockGroupRepository struct {
	groups []models.Group
}

func (m *MockGroupRepository) CreateGroup(groupName *string, userIds []primitive.ObjectID) (*models.Group, error) {
	group := &models.Group{
		DefaultModel: mgm.DefaultModel{
			IDField:    mgm.IDField{ID: primitive.NewObjectID()},
			DateFields: mgm.DateFields{},
		},
		Name:    groupName,
		UserIds: userIds,
	}

	m.groups = append(m.groups, *group)

	return group, nil
}

func (m MockGroupRepository) GetGroups(userId primitive.ObjectID) ([]models.Group, error) {
	var filteredGroups []models.Group
	for _, group := range m.groups {
		for _, groupUserId := range group.UserIds {
			if groupUserId == userId {
				filteredGroups = append(filteredGroups, group)
				break
			}
		}
	}

	return filteredGroups, nil
}

func (m MockGroupRepository) GroupExists(userIds []primitive.ObjectID) (bool, error) {
	for _, group := range m.groups {
		if comparePrimitiveIdsArray(group.UserIds, userIds) {
			return true, nil
		}
	}

	return false, nil
}

func comparePrimitiveIdsArray(x, y []primitive.ObjectID) bool {
	if len(x) != len(y) {
		return false
	}
	diff := make(map[primitive.ObjectID]int, len(x))
	for _, _x := range x {
		diff[_x]++
	}
	for _, _y := range y {
		if _, ok := diff[_y]; !ok {
			return false
		}
		diff[_y] -= 1
		if diff[_y] == 0 {
			delete(diff, _y)
		}
	}
	if len(diff) == 0 {
		return true
	}
	return false
}
