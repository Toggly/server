package models

import "time"

// ProjectStatus enum
const (
	ProjectStatusActive   = "active"
	ProjectStatusDisabled = "disabled"
)

// Project type
type Project struct {
	Code        string    `json:"code"`
	Owner       string    `json:"owner"`
	Status      string    `json:"status"`
	Description string    `json:"description"`
	RegDate     time.Time `json:"reg_date" bson:"reg_date"`
}
