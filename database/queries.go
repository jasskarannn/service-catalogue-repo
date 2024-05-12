package database

import "database/sql"

// vars for Service
var (
	queryInsertToService          = `INSERT INTO service (name, description, version_count, created_at) VALUES ($1, $2, $3, $4)`
	queryGetServices              = `SELECT service_id, name, description, version_count, coalesce(created_at, $1) FROM service OFFSET $2 LIMIT $3`
	querySearchServices           = `SELECT service_id, name, description, version_count, COALESCE(created_at, $1) FROM service WHERE name ILIKE '%' || $2 || '%'`
	queryGetVersionCountOfService = `SELECT COUNT(*) AS version_count FROM version WHERE service_id = $1`
	queryUpdateToService          = `UPDATE service SET version_count = $1 WHERE service_id = $2`
	queryGetServiceByID           = `SELECT service_id, name, description, version_count, created_at FROM service WHERE service_id = $1`
)

// vars for Version
var (
	queryAddVersion             = `INSERT INTO version (service_id, version_number, description, created_at) VALUES ($1, $2, $3, $4)`
	queryServiceInfoFromVersion = `SELECT version_id, service_id, version_number, description, created_at FROM version WHERE service_id = $1`
)

func GetVersionCountByServiceID(tx *sql.Tx, serviceID int) (int, error) {
	var versionCount int
	err := tx.QueryRow(queryGetVersionCountOfService, serviceID).Scan(&versionCount)
	if err != nil {
		return 0, err
	}
	return versionCount, nil
}
