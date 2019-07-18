package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Environment type
type Environment struct {
	ID          primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	ProjectID   primitive.ObjectID `json:"-" bson:"project_id"`
	Code        string             `json:"code"`
	Description string             `json:"description"`
	Protected   bool               `json:"protected"`
}

// EnvAPIKey struct
type EnvAPIKey struct {
	ID          primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	EnvID       primitive.ObjectID `json:"-" bson:"env_id"`
	Name        string             `json:"name"`
	Key         string             `json:"key"`
	Secret      string             `json:"secret"`
	CreatedDate time.Time          `json:"createdDate" bson:"create_date"`
	UsedDate    time.Time          `json:"usedDate" bson:"used_date"`
}
