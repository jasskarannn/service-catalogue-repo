package database

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jasskarannn/service-catalogue/config"
	_ "github.com/lib/pq"
)

// InitDB initializes the database connection
func InitDB(cfg config.Config) (*sql.DB, error) {
	// Connection string for PostgreSQL database
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.DBName)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}
	// Check the connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

// GenerateUUID generates a new UUID
func GenerateUUID() string {
	return uuid.New().String()
}
