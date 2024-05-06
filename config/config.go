package config

import (
	"encoding/json"
	"os"
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
func LoadConfig(filePath string) (Config, error) {
	var cfg Config
	file, err := os.Open(filePath)
	if err != nil {
		return cfg, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&cfg)
	return cfg, err
}
