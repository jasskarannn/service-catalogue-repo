package database

import (
	"database/sql"

	"github.com/jasskarannn/service-catalogue/helpers"
	"github.com/jasskarannn/service-catalogue/models"
)

type VersionRepository interface {
	AddVersion(version models.Version) error
	GetVersionsByServiceID(serviceID string) ([]models.Version, error)
}

type versionRepositoryImpl struct {
	db *sql.DB
}

// NewVersionRepository creates a new instance of the versionRepositoryImpl.
func NewVersionRepository(db *sql.DB) VersionRepository {
	return &versionRepositoryImpl{db: db}
}

// AddVersion adds a new version to the database
func (r *versionRepositoryImpl) AddVersion(version models.Version) error {
	// Begin transaction
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	// Rollback transaction if any error occurs
	defer tx.Rollback()

	// Execute the SQL query to add the version
	hashedUUID := helpers.HashedUUIDAsTimestamp()
	_, err = tx.Exec(queryAddVersion, version.ServiceID, version.VersionNumber, version.Description, hashedUUID)
	if err != nil {
		return err
	}

	// Get the current version count for the service
	versionCount, err := GetVersionCountByServiceID(tx, version.ServiceID)
	if err != nil {
		return err
	}

	// Increment the version count
	versionCount++

	// Update the version count in the service table
	_, err = tx.Exec(queryUpdateToService, versionCount, version.ServiceID)
	if err != nil {
		return err
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// GetVersionsByServiceID fetches the versions associated with a service from the database
func (r *versionRepositoryImpl) GetVersionsByServiceID(serviceID string) ([]models.Version, error) {
    // Execute the SQL query
    rows, err := r.db.Query(queryServiceInfoFromVersion, serviceID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    // Initialize a slice to store the fetched versions
    var versions []models.Version

    // Iterate over the rows and populate the versions slice
    for rows.Next() {
        var version models.Version
        // Scan the row data into the Version struct
        err := rows.Scan(&version.VersionID, &version.ServiceID, &version.VersionNumber, &version.Description, &version.CreatedAt)
        if err != nil {
            return nil, err
        }
        versions = append(versions, version)
    }

    // Check for any errors during iteration
    if err := rows.Err(); err != nil {
        return nil, err
    }

    return versions, nil
}
