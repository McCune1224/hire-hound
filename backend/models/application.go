package models

import (
	"time"

	"gorm.io/gorm"
)

type Application struct {
	gorm.Model
	JobTitle    string    `json:"jobTitle"`
	Location    string    `json:"location"`
	Description string    `json:"description"`
	DateApplied time.Time `json:"dateApplied"`
	Status      string    `json:"status"`
	// Company foreign key
	CompanyID uint `json:"companyID"`
	// User foreign key
	UserID uint `json:"userID"`
}
