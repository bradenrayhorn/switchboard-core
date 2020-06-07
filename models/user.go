package models

import (
	"github.com/Kamva/mgm/v3"
)

type User struct {
	mgm.DefaultModel `bson:",inline"`
	Username         string
	Password         string
}
