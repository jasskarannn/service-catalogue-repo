package database

import (
	"database/sql"
	"fmt"
)

// ServiceRepository defines the interface for interacting with the service table in the database.
type ServiceRepository interface {
	GetServices(offset, limit int) ([]Service, error)
}

// serviceRepositoryImpl is an implementation of the ServiceRepository interface.
type serviceRepositoryImpl struct {
	db *sql.DB
}

// NewServiceRepository creates a new instance of the serviceRepositoryImpl.
func NewServiceRepository(db *sql.DB) ServiceRepository {
	return &serviceRepositoryImpl{db: db}
}

// GetServices retrieves services from the database with pagination.
func (r *serviceRepositoryImpl) GetServices(offset, limit int) ([]Service, error) {
	query := `SELECT * FROM public."service" OFFSET $1 LIMIT $2`
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	fmt.Println("<<< rows is ", rows)

	var services []Service
	for rows.Next() {
		var service Service
		err := rows.Scan(&service.ServiceID, &service.Name, &service.Description, &service.Version, &service.CreatedAt, &service.UpdatedAt)
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
