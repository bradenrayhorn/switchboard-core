package models

import (
	"github.com/Kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Group struct {
	mgm.DefaultModel `bson:",inline"`
	Name             *string              `json:"name"`
	UserIds          []primitive.ObjectID `json:"users" bson:"users"`
}

type GroupUser struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type GroupResponse struct {
	ID    primitive.ObjectID `json:"id"`
	Name  *string            `json:"name"`
	Users []GroupUser        `json:"users"`
}
