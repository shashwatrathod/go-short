package handlers

import (
	"encoding/json"
	"net/http"
)

type HealthResponse struct {
	Status string `json:"status"`
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	// Create a response object
	response := HealthResponse{
		Status: "ok",
	}

	// Encode the response as JSON and write it to the response writer
	json.NewEncoder(w).Encode(response)
}