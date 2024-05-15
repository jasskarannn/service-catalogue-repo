package database

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrServiceNotFound = errors.New("service not found")
)

// Service represents a service in the catalog.
type Service struct {
	ServiceID    int       `json:"service_id" db:"service_id"`
	Name         string    `json:"name" db:"name"`
	Description  string    `json:"description" db:"description"`
	VersionCount string    `json:"version_count" db:"version_count"`
	CreatedAt    uuid.UUID `json:"created_at" db:"created_at"`
}
