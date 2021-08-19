package model

import (
	"github.com/google/uuid"
)

// User represents a single user that can use API
type User struct {
	ID       uuid.UUID `bson:"_id"`
	Login    string
	Password string
	Admin    bool
}

// UserRequest struct is used for binding the request content
type UserRequest struct {
	Login    string `query:"login" form:"login" json:"login" validate:"alphanum,required"`
	Password string `query:"password" form:"password" json:"password" validate:"required,min=4"`
	Admin    bool   `query:"admin" form:"admin" json:"admin"`
}
