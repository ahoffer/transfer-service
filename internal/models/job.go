package models

import (
	"time"

	"gorm.io/gorm"
)

// JobRequest represents a job transfer request
type JobRequest struct {
	gorm.Model
	Name           string    `json:"name" gorm:"not null"`
	SourceUrl      string    `json:"sourceUrl" gorm:"not null"`
	Destination    string    `json:"destination" gorm:"not null"`
	DestinationUrl string    `json:"destinationUrl" gorm:"not null"`
	CreatedAt      time.Time `json:"createdAt" gorm:"not null"`
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
