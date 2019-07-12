package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Object struct
type Object struct {
	ID         primitive.ObjectID     `json:"-" bson:"_id"`
	InstanceID string                 `json:"instanceId" bson:"instance_id"`
	Name       string                 `json:"name"`
	Props      map[string]interface{} `json:"props"`
	ProjectID  primitive.ObjectID     `json:"projectId" bson:"project_id"`
	OwnerID    primitive.ObjectID     `json:"ownerId" bson:"owner_id"`
}
