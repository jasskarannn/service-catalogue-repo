package database

import "github.com/google/uuid"

// Version represents a version in the catalog.
type Version struct {
	VersionID     int       `json:"version_id" db:"version_id"`
	ServiceID     int       `json:"service_id" db:"service_id"`
	VersionNumber int       `json:"version_number" db:"version_number"`
	Description   string    `json:"description" db:"description"`
	CreatedAt     uuid.UUID `json:"created_at" db:"created_at"`
}
