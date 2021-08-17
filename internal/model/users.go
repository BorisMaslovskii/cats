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
