package config

import (
	"log"

	"github.com/go-ini/ini"
)

// Config struct to hold configuration data
type Config struct {
	Database struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
		DBName   string `json:"dbname"`
	} `json:"database"`
	Server struct {
		Port int `json:"port"`
	} `json:"server"`
}

// LoadConfig loads configuration from a JSON file
func LoadConfig(filePath string) *ini.File {
	cfg, err := ini.Load(filePath)
	if err != nil {
		log.Fatalf("failed to load configuration file: %v", err)
	}
	return cfg
}
