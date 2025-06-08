package handlers

import (
	"encoding/json"
	"net/http"
)

// HealthResponse defines the response for the health check endpoint.
// @Description Response for the health check endpoint.
type HealthResponse struct {
	Status string `json:"status" example:"ok"`
}

// HealthHandler serves the health check endpoint.
//
// @Summary Application health check
// @Description Returns the health status of the application.
// @Tags health
// @Produce json
// @Success 200 {object} HealthResponse "Application is healthy"
// @Router /health [get]
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	// Create a response object
	response := HealthResponse{
		Status: "ok",
	}
	w.Header().Set("Content-Type", "application/json")
	// Encode the response as JSON and write it to the response writer
	json.NewEncoder(w).Encode(response)
}
