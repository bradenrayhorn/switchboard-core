package repositories

import (
	"github.com/Kamva/mgm/v3"
	"github.com/Kamva/mgm/v3/operator"
	"github.com/bradenrayhorn/switchboard-core/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var User UserRepository

func init() {
	User = MongoUserRepository{}
}

type UserRepository interface {
	CreateUser(username string, hashedPassword string) (*models.User, error)
	GetUser(username string) (*models.User, error)
	Exists(username string) (bool, error)
	GetUsers(userIDs []primitive.ObjectID) ([]models.User, error)
}

type MongoUserRepository struct{}

func (r MongoUserRepository) CreateUser(username string, hashedPassword string) (*models.User, error) {
	user := &models.User{
		Username: username,
		Password: hashedPassword,
	}

	err := mgm.Coll(user).Create(user)
	return user, err
}

func (r MongoUserRepository) GetUser(username string) (*models.User, error) {
	user := &models.User{}

	err := mgm.Coll(user).First(bson.M{"username": username}, user)

	return user, err
}

func (r MongoUserRepository) Exists(username string) (bool, error) {
	cursor, err := mgm.Coll(&models.User{}).Find(mgm.Ctx(), bson.M{"username": username})
	if err != nil {
		return false, err
	}
	cursor.Next(mgm.Ctx())

	return cursor.Current != nil, cursor.Close(mgm.Ctx())
}

func (r MongoUserRepository) GetUsers(userIDs []primitive.ObjectID) ([]models.User, error) {
	var users = make([]models.User, 0)
	cursor, err := mgm.Coll(&models.User{}).Find(mgm.Ctx(), bson.M{"_id": bson.M{operator.In: userIDs}})
	if err != nil {
		return users, err
	}

	err = cursor.All(mgm.Ctx(), &users)

	return users, nil
}
