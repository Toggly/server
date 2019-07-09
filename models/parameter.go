package models

// Parameter types enum
const (
	ParameterTypeBool   = "bool"
	ParameterTypeString = "string"
	ParameterTypeInt    = "int"
)

// Parameter type
type Parameter struct {
	Code string `json:"code"`
	// Project       string        `json:"project"`
	// Environment   string        `json:"environment"`
	// Group         string        `json:"group"`
	Description   string        `json:"description"`
	Type          string        `json:"type"`
	Value         interface{}   `json:"value"`
	AllowedValues []interface{} `json:"allowed_values,omitempty" bson:"allowed_values,omitempty"`
}
