package database

import (
	"time"
)

// Service represents a service in the catalog.
type Service struct {
	ServiceID   int       `json:"service_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Version     string    `json:"version"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
