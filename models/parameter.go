package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Parameter types enum
const (
	ParameterTypeBool   = "bool"
	ParameterTypeString = "string"
	ParameterTypeInt    = "int"
)

// Parameter type
type Parameter struct {
	ID            primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	ProjectID     primitive.ObjectID `json:"-" bson:"project_id"`
	Code          string             `json:"code"`
	Description   string             `json:"description"`
	Type          string             `json:"type"`
	Value         interface{}        `json:"value"`
	AllowedValues []interface{}      `json:"allowed_values,omitempty" bson:"allowed_values,omitempty"`
}
