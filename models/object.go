package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Object struct
type Object struct {
	ID         primitive.ObjectID     `json:"-" bson:"_id,omitempty"`
	InstanceID string                 `json:"instanceId" bson:"instance_id"`
	Name       string                 `json:"name"`
	Props      map[string]interface{} `json:"props"`
	EnvID      primitive.ObjectID     `json:"envId" bson:"env_id"`
	ProjectID  primitive.ObjectID     `json:"projectId" bson:"project_id"`
}
