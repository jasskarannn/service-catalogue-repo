package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/jasskarannn/service-catalogue/database"
	"net/http"
	"strconv"
	"sync"
)

// Handler struct holds the dependencies for HTTP handlers
type Handler struct {
	DB                *sql.DB
	ServiceRepository database.ServiceRepository
	mu                sync.RWMutex
}

// NewHandler creates a new instance of Handler
func NewHandler(db *sql.DB, serviceRepo database.ServiceRepository) *Handler {
	return &Handler{
		DB:                db,
		ServiceRepository: serviceRepo,
	}
}

// SetupRouter sets up the HTTP router and registers handlers
func SetupRouter(h *Handler) http.Handler {
	r := http.NewServeMux()
	r.HandleFunc("/db_health", h.CheckDBHealthHandler)
	r.HandleFunc("/services", h.GetServicesHandler)
	r.HandleFunc("/services/search", h.SearchServicesHandler)
	//r.HandleFunc("/services", h.addServiceHandler)
	return r
}

func (h *Handler) CheckDBHealthHandler(w http.ResponseWriter, r *http.Request) {
	// Perform a simple query to check the database connection
	_, err := h.DB.Query("SELECT 1")
	if err != nil {
		// Database connection error
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(w, "Database connection error: %v", err)
		return
	}

	// Database is healthy
	healthResponse := struct {
		Status string `json:"status"`
	}{
		Status: "ok",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(healthResponse)
	if err != nil {
		return
	}
}
func (h *Handler) GetServicesHandler(w http.ResponseWriter, r *http.Request) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	// Parse query parameters for pagination
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))   // Page number (default: 1)
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit")) // Number of items per page (default: 10)

	// Implement pagination logic
	// Calculate offset based on page and limit
	offset := (page - 1) * limit

	//Retrieve services from the database
	services, err := h.ServiceRepository.GetServices(offset, limit)
	if err != nil {
		fmt.Println("[getServicesHandler] error ", err)
	}
	//Replace the above line with your actual database retrieval logic

	// Simulated example data for demonstration
	//services := []string{"Service 1", "Service 2", "Service 3"}

	// Slice services based on pagination parameters
	startIndex := offset
	endIndex := offset + limit
	if endIndex > len(services) {
		endIndex = len(services)
	}
	paginatedServices := services[startIndex:endIndex]

	// Prepare JSON response
	response := map[string]interface{}{
		"page":     page,
		"limit":    limit,
		"services": paginatedServices,
	}

	// Write JSON response to the response writer
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) SearchServicesHandler(w http.ResponseWriter, r *http.Request) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	// Retrieve search query from URL parameters
	// Search services in the database based on the query
	// Return JSON response
}
