package models

import (
	"github.com/Kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRole string

const (
	RoleUser  UserRole = "user"
	RoleAdmin          = "admin"
)

type Organization struct {
	mgm.DefaultModel `bson:",inline"`
	Name             string             `json:"name"`
	Users            []OrganizationUser `json:"users" bson:"users"`
}

type OrganizationUser struct {
	ID   primitive.ObjectID
	Role UserRole
}
