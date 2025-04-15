package models

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Employee struct {
	ID         bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name       string        `json:"name,omitempty" bson:"name,omitempty" validate:"required"`
	Email      string        `json:"email,omitempty" bson:"email,omitempty" validate:"required,email"`
	Phone      string        `json:"phone,omitempty" bson:"phone,omitempty" validate:"required"`
	Department string        `json:"department,omitempty" bson:"department,omitempty" validate:"required"`
}
