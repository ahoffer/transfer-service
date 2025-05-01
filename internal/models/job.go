package models

import (
	"time"
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
	ID        string    `json:"id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	JobRequest
}
