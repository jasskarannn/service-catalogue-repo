package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/jasskarannn/service-catalogue/database"
	"github.com/jasskarannn/service-catalogue/models"
)

// Handler struct holds the dependencies for HTTP handlers
type Handler struct {
	DB                *sql.DB
	ServiceRepository database.ServiceRepository
	VersionRepository database.VersionRepository
	mu                sync.RWMutex
}

// NewHandler creates a new instance of Handler
func NewHandler(db *sql.DB, serviceRepo database.ServiceRepository, versionRepo database.VersionRepository) *Handler {
	return &Handler{
		DB:                db,
		ServiceRepository: serviceRepo,
		VersionRepository: versionRepo,
	}
}

// SetupRouter sets up the HTTP router and registers handlers
func SetupRouter(h *Handler) http.Handler {
	r := http.NewServeMux()
	r.HandleFunc("/db_health", h.CheckDBHealthHandler)
	r.HandleFunc("/services", h.GetServicesHandler)
	r.HandleFunc("/services/search", h.SearchServicesHandler)
	r.HandleFunc("/add_service", h.AddServiceHandler)
	r.HandleFunc("/add_version", h.AddVersionHandler)
	r.HandleFunc("/all_service_details", h.GetServiceWithVersionsHandler)
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

// AddServiceHandler handles the POST request to add a new service to the database
func (h *Handler) AddServiceHandler(w http.ResponseWriter, r *http.Request) {
	// Decode the request body into a Service object
	var service models.Service
	err := json.NewDecoder(r.Body).Decode(&service)
	fmt.Println("service : ", service)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Insert the new service into the database
	err = h.ServiceRepository.AddService(service)
	if err != nil {
		http.Error(w, "Failed to add service to database", http.StatusInternalServerError)
		fmt.Println("[AddServiceHandler] err : ", err)
		return
	}

	// Respond with success message
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Service added successfully"))
}

func (h *Handler) GetServicesHandler(w http.ResponseWriter, r *http.Request) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	// Parse query parameters for pagination
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))   // Page number (default: 1)
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit")) // Number of items per page (default: 10)
	serviceName := r.URL.Query().Get("service_name")     // If passed, will cater the usecase to show a specific service card

	// Implement pagination logic
	// Calculate offset based on page and limit
	offset := (page - 1) * limit
	if offset < 0 {
		offset = 0
		page = 1
	}

	//Retrieve services from the database
	services, err := h.ServiceRepository.GetServices(offset, limit, serviceName)
	if err != nil {
		fmt.Println("[getServicesHandler] error ", err)
	}

	// Prepare JSON response
	response := map[string]interface{}{
		"page":     page,
		"limit":    limit,
		"services": services,
	}

	// Write JSON response to the response writer
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) SearchServicesHandler(w http.ResponseWriter, r *http.Request) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	// Retrieve search query from URL parameters
	query := r.URL.Query().Get("query")

	// Search services in the database based on the query
	services, err := h.ServiceRepository.SearchServices(query)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	response := map[string]interface{}{
		"services": services,
	}

	// Write JSON response to the response writer
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) AddVersionHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the request body to get the version information
	var version models.Version
	err := json.NewDecoder(r.Body).Decode(&version)
	if err != nil {
		fmt.Println("error 1 : ", err)
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Validate the version information
	if version.ServiceID == 0 || version.VersionNumber == 0.0 {
		http.Error(w, "Invalid version data", http.StatusBadRequest)
		return
	}

	// Add the version to the database
	err = h.VersionRepository.AddVersion(version)
	if err != nil {
		http.Error(w, "Failed to add version", http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Version added successfully"))
}

func (h *Handler) GetServiceWithVersionsHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the service_id from the request parameters
	serviceID := r.URL.Query().Get("service_id")

	// Query the service table to retrieve the service details
	service, err := h.ServiceRepository.GetServiceByID(serviceID)
	if err != nil {
		// Handle error (e.g., service not found)
		http.Error(w, "Failed to retrieve service details", http.StatusInternalServerError)
		return
	}

	// Query the version table to retrieve the related versions for the service
	versions, err := h.VersionRepository.GetVersionsByServiceID(serviceID)
	if err != nil {
		// Handle error (e.g., versions not found)
		http.Error(w, "Failed to retrieve versions", http.StatusInternalServerError)
		return
	}

	// Assemble the response JSON object
	response := struct {
		Service  models.Service   `json:"service"`
		Versions []models.Version `json:"versions"`
	}{
		Service:  *service,
		Versions: versions,
	}

	// Marshal the response to JSON
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		// Handle error (e.g., JSON marshaling error)
		http.Error(w, "Failed to marshal JSON response", http.StatusInternalServerError)
		return
	}

	// Set the content type header and write the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
