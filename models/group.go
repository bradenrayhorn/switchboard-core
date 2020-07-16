package models

import (
	"github.com/Kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GroupType string

const (
	GroupTypePublicChannel  GroupType = "public_channel"
	GroupTypePrivateChannel           = "private_channel"
	GroupTypePrivateMessage           = "private_message"
)

type Group struct {
	mgm.DefaultModel `bson:",inline"`
	Name             *string            `json:"name"`
	Type             GroupType          `json:"type"`
	Users            []GroupUser        `json:"users" bson:"users"`
	Organization     primitive.ObjectID `json:"organization" bson:"organization"`
}

type GroupUser struct {
	ID       primitive.ObjectID `json:"id"`
	Username string             `json:"username"`
}

type GroupResponse struct {
	ID           primitive.ObjectID `json:"id"`
	Name         *string            `json:"name"`
	Users        []GroupUser        `json:"users"`
	Organization primitive.ObjectID `json:"organization"`
}
