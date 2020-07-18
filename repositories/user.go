package repositories

import (
	"github.com/Kamva/mgm/v3"
	"github.com/bradenrayhorn/switchboard-core/models"
	"go.mongodb.org/mongo-driver/bson"
)

var User UserRepository

func init() {
	User = MongoUserRepository{}
}

type UserRepository interface {
	CreateUser(username string, hashedPassword string) (*models.User, error)
	GetUser(username string) (*models.User, error)
	Exists(username string) (bool, error)
	DropAll() error
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

func (r MongoUserRepository) DropAll() error {
	return mgm.Coll(&models.User{}).Drop(mgm.Ctx())
}
