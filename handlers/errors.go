package handlers

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse is the structure for a 404 error response.
type ErrorResponse struct {
    Error   string `json:"error"`
    Message string `json:"message"`
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	// Set the response status to 404 Not Found
	w.WriteHeader(http.StatusNotFound)

	// Create a NotFoundResponse object
	response := ErrorResponse{
		Error:   "Not Found",
		Message: "The requested resource was not found.",
	}

	// Encode the response as JSON and write it to the response writer
	json.NewEncoder(w).Encode(response)
}

func MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	// Set the response status to 405 Method Not Allowed
	w.WriteHeader(http.StatusMethodNotAllowed)

	// Create a MethodNotAllowedResponse object
	response := ErrorResponse{
		Error:   "Method Not Allowed",
		Message: "The requested method is not allowed for this resource.",
	}

	// Encode the response as JSON and write it to the response writer
	json.NewEncoder(w).Encode(response)
}