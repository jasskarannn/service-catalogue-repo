package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jasskarannn/service-catalogue/helpers"
	"github.com/jasskarannn/service-catalogue/models"
)

// ServiceRepository defines the interface for interacting with the service table in the database.
type ServiceRepository interface {
	AddService(service models.Service) error
	GetServices(offset, limit int, serviceName string) ([]Service, error)
	SearchServices(query string) ([]Service, error)
	GetServiceByID(serviceID string) (*models.Service, error)
}

// serviceRepositoryImpl is an implementation of the ServiceRepository interface.
type serviceRepositoryImpl struct {
	db *sql.DB
}

// NewServiceRepository creates a new instance of the serviceRepositoryImpl.
func NewServiceRepository(db *sql.DB) ServiceRepository {
	return &serviceRepositoryImpl{db: db}
}

// AddService adds a new service to the database
func (r *serviceRepositoryImpl) AddService(service models.Service) error {
	// Begin transaction
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	// Rollback transaction if any error occurs
	defer tx.Rollback()

	hashedUUID := helpers.HashedUUIDAsTimestamp()
	versionCount, err := GetVersionCountByServiceID(tx, service.ServiceID)
	if err != nil {
		fmt.Println("[AddService] error ", err)
	}

	if versionCount > 0 {
		// Increment the version count in the existing service
		versionCount++
		_, err = tx.Exec(queryUpdateToService, versionCount, service.ServiceID)
		if err != nil {
			return err
		}
	} else {
		// Add a new service
		_, err = tx.Exec(queryInsertToService, service.Name, service.Description, versionCount, hashedUUID)
		if err != nil {
			return err
		}
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// GetServices retrieves services from the database with pagination.
func (r *serviceRepositoryImpl) GetServices(offset, limit int, serviceName string) ([]Service, error) {
	rows, err := r.db.Query(queryGetServices, uuid.Nil, offset, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	fmt.Println("<<< rows is ", rows)

	var services []Service
	for rows.Next() {
		var service Service
		err := rows.Scan(&service.ServiceID, &service.Name, &service.Description, &service.VersionCount, &service.CreatedAt)
		if err != nil {
			return nil, err
		}
		services = append(services, service)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return services, nil
}

// SearchServices searches for services in the database based on the query string
func (r *serviceRepositoryImpl) SearchServices(query string) ([]Service, error) {
	rows, err := r.db.Query(querySearchServices, uuid.Nil, query)
	if err != nil {
		fmt.Println("[SearchServices] error ", err)
		return nil, err
	}
	defer rows.Close()

	// Iterate over the rows and populate the services slice
	var services []Service
	for rows.Next() {
		var service Service
		if err := rows.Scan(&service.ServiceID, &service.Name, &service.Description, &service.VersionCount, &service.CreatedAt); err != nil {
			return nil, err
		}
		services = append(services, service)
	}

	// Check for any errors during iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return services, nil
}

// GetServiceByID fetches the service details by service ID from the database
func (r *serviceRepositoryImpl) GetServiceByID(serviceID string) (*models.Service, error) {
	// Execute the SQL query
	row := r.db.QueryRow(queryGetServiceByID, serviceID)

	// Initialize a new Service struct to store the fetched data
	var service models.Service

	// Scan the row data into the Service struct
	err := row.Scan(&service.ServiceID, &service.Name, &service.Description, &service.VersionCount, &service.CreatedAt)
	if err != nil {
		// Check if no rows were returned (i.e., service with the given ID does not exist)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrServiceNotFound
		}
		return nil, err
	}

	return &service, nil
}
