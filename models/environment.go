package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Environment type
type Environment struct {
	ID        primitive.ObjectID `json:"-" bson:"_id"`
	Code      string             `json:"code"`
	ProjectID primitive.ObjectID `json:"projectId" bson:"project_id"`
	OwnerID   primitive.ObjectID `json:"ownerId" bson:"owner_id"`
	Protected bool               `json:"protected"`
	RegDate   time.Time          `json:"reg_date" bson:"reg_date"`
}
