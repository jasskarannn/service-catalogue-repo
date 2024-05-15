package database

import (
	"database/sql"
	"fmt"

	"github.com/go-ini/ini"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

// InitDB initializes the database connection
func InitDB(cfg *ini.File) (*sql.DB, error) {
	// Connection string for PostgreSQL database
	dbSection := cfg.Section("database")
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbSection.Key("host").String(),
		dbSection.Key("port").MustInt(),
		dbSection.Key("user").String(),
		dbSection.Key("password").String(),
		dbSection.Key("dbname").String())

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
