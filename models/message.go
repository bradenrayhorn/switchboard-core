package models

import (
	"github.com/Kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	mgm.DefaultModel `bson:",inline"`
	SenderUserId     primitive.ObjectID
	GroupId          primitive.ObjectID
	Message          string
}
