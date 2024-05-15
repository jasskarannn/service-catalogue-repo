package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jasskarannn/service-catalogue/config"
	"github.com/jasskarannn/service-catalogue/database"
	"github.com/jasskarannn/service-catalogue/handlers"
)

func main() {
	// Load Configuration based on app environment
	var configFile string
	env := "development"
	if env == "development" {
		configFile = "config/dev-config.ini"
	} else {
		configFile = "config/prod-config.ini"
	}

	cfg := config.LoadConfig(configFile)

	// Initialize the database
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("error initializing database: %s", err)
	}

	// Initialize the service repository
	serviceRepo := database.NewServiceRepository(db)

	// Initialize the version repository
	versionRepo := database.NewVersionRepository(db)

	// Initialize handlers
	h := handlers.NewHandler(db, serviceRepo, versionRepo)

	// Initialize router
	r := handlers.SetupRouter(h)

	// Start the server
	serverSection := cfg.Section("server")
	addr := fmt.Sprintf(":%d", serverSection.Key("port").MustInt())
	log.Printf("app server listening on port %d", serverSection.Key("port").MustInt())
	log.Fatal(http.ListenAndServe(addr, r))
}
