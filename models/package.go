package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Package struct
type Package struct {
	ID        primitive.ObjectID `json:"-" bson:"_id"`
	Code      string             `json:"code"`
	Name      string             `json:"name"`
	ProjectID primitive.ObjectID `json:"projectId" bson:"project_id"`
	OwnerID   primitive.ObjectID `json:"ownerId" bson:"owner_id"`
}
