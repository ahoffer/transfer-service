package models

import (
	"gorm.io/gorm"
)

// JobRequest represents the incoming request to create a new transfer job
type JobRequest struct {
	Name           string `json:"name"`
	SourceUrl      string `json:"sourceUrl"`
	Destination    string `json:"destination"`
	DestinationUrl string `json:"destinationUrl"`
}

// Job represents a transfer job in the system
type Job struct {
	gorm.Model
	Name           string `json:"name" gorm:"not null"`
	SourceUrl      string `json:"sourceUrl" gorm:"not null"`
	Destination    string `json:"destination" gorm:"not null"`
	DestinationUrl string `json:"destinationUrl" gorm:"not null"`
	Status         string `json:"status" gorm:"not null;default:'created'"`
}
