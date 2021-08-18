// Package model provides structs required for a business logic
package model

import (
	"github.com/google/uuid"
)

// Cat represents a single cat
type Cat struct {
	ID    uuid.UUID `bson:"_id"`
	Name  string
	Color string
}

// CatRequest struct is used for binding the request content
type CatRequest struct {
	Name  string `form:"name" json:"name" validate:"alphanum,required"`
	Color string `form:"color" json:"color" validate:"alpha"`
}
