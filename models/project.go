package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ProjectStatus enum
const (
	ProjectStatusActive   = "active"
	ProjectStatusDisabled = "disabled"
)

// Project type
type Project struct {
	ID          primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	Code        string             `json:"code"`
	Name        string             `json:"name"`
	OwnerID     string             `json:"-" bson:"owner_id"`
	Status      int                `json:"status"`
	Description string             `json:"description"`
	RegDate     time.Time          `json:"reg_date" bson:"reg_date"`
}
