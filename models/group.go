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
