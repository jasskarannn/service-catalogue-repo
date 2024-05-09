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
	// Load configuration
	cfg, err := config.LoadConfig("config/config.json")
	if err != nil {
		log.Fatalf("Error loading config: %s", err)
	}

	// Initialize the database
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("Error initializing database: %s", err)
	}

	// Initialize the service repository
	serviceRepo := database.NewServiceRepository(db)

	// Initialize handlers
	h := handlers.NewHandler(db, serviceRepo)

	// Initialize router
	r := handlers.SetupRouter(h)

	// Start the server
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("Server listening on port %d", cfg.Server.Port)
	log.Fatal(http.ListenAndServe(addr, r))
}
