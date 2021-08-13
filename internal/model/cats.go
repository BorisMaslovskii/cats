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

// User represents a single user that can use API
type User struct {
	ID       uuid.UUID `bson:"_id"`
	Login    string
	Password string
}
