package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Package struct
type Package struct {
	ID          primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	ProjectID   primitive.ObjectID `json:"-" bson:"project_id"`
	Code        string             `json:"code"`
	Description string             `json:"description"`
}
