package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Parameter types enum
const (
	ParameterTypeBool   = "bool"
	ParameterTypeString = "string"
	ParameterTypeInt    = "int"
	ParameterTypeFloat  = "float64"
	ParameterTypeEnum   = "enum"
)

// Parameter type
type Parameter struct {
	ID          primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	ProjectID   primitive.ObjectID `json:"-" bson:"project_id"`
	Code        string             `json:"code"`
	Description string             `json:"description"`
	Type        string             `json:"type"`
	Value       interface{}        `json:"value"`
	Updated     time.Time          `json:"updated"`
}

// ParameterValue final struct
type ParameterValue struct {
	Code  string      `json:"code"`
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}
