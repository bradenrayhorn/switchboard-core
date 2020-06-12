package repositories

import (
	"errors"
	"github.com/Kamva/mgm/v3"
	"github.com/bradenrayhorn/switchboard-backend/models"
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

// mock repository

type MockUserRepository struct {
	users []models.User
}

func (r *MockUserRepository) CreateUser(username string, hashedPassword string) (*models.User, error) {
	user := &models.User{
		DefaultModel: mgm.DefaultModel{
			IDField:    mgm.IDField{ID: primitive.NewObjectID()},
			DateFields: mgm.DateFields{},
		},
		Username: username,
		Password: hashedPassword,
	}

	r.users = append(r.users, *user)
	return user, nil
}

func (r MockUserRepository) GetUser(username string) (*models.User, error) {
	for i := range r.users {
		if r.users[i].Username == username {
			return &r.users[i], nil
		}
	}

	return nil, errors.New("no user found")
}

func (r MockUserRepository) Exists(username string) (bool, error) {
	for i := range r.users {
		if r.users[i].Username == username {
			return true, nil
		}
	}

	return false, nil
}
