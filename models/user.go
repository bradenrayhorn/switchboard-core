package models

import (
	"github.com/Kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	mgm.DefaultModel `bson:",inline"`
	Username         string
	Password         string
}

type UserResponse struct {
	ID       primitive.ObjectID `json:"id"`
	Username string             `json:"username"`
}
