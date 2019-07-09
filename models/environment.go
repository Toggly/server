package models

import "time"

// Environment type
type Environment struct {
	Code      string    `json:"code"`
	Project   string    `json:"project"`
	Protected bool      `json:"protected"`
	RegDate   time.Time `json:"reg_date" bson:"reg_date"`
}
