package repositories

import (
	"github.com/Kamva/mgm/v3"
	"github.com/bradenrayhorn/switchboard-backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(username string, password string) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username: username,
		Password: string(hashedPassword),
	}

	err = mgm.Coll(user).Create(user)
	return user, err
}

func GetUser(username string) (*models.User, error) {
	user := &models.User{}

	err := mgm.Coll(user).First(bson.M{"username": username}, user)

	return user, err
}

func UserExists(username string) (bool, error) {
	cursor, err := mgm.Coll(&models.User{}).Find(mgm.Ctx(), bson.M{"username": username})
	if err != nil {
		return false, err
	}
	if cursor.Current == nil {
		return false, nil
	}
	return true, cursor.Close(mgm.Ctx())
}
