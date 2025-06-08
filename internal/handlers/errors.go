package handlers

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse is the structure for a generic error response.
// @Description Generic error response structure used for 4xx and 5xx errors.
type ErrorResponse struct {
	Error   string `json:"error" example:"Error Type"`                     // Type of the error (e.g., "Not Found", "Internal Server Error")
	Message string `json:"message" example:"A descriptive error message."` // Detailed error message
}

// NotFoundHandler handles requests for routes that are not found.
// It's typically used as a global NotFound handler for the router.
//
// @Summary Not Found
// @Description Handles requests for routes that are not found.
// @Tags errors
// @Produce json
// @Success 404 {object} ErrorResponse "Resource not found"
// @Router /anyNonExistentRoute [get]
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	// Set the response status to 404 Not Found
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)

	// Create a NotFoundResponse object
	response := ErrorResponse{
		Error:   "Not Found",
		Message: "The requested resource was not found.",
	}

	// Encode the response as JSON and write it to the response writer
	json.NewEncoder(w).Encode(response)
}

// MethodNotAllowedHandler handles requests where the HTTP method is not allowed for the route.
// It's typically used as a global MethodNotAllowed handler for the router.
// @Summary Method Not Allowed
// @Description Handles requests where the HTTP method is not allowed for the route.
// @Tags errors
// @Produce json
// @Success 405 {object} ErrorResponse "Method not allowed"
// @Router /anyRouteWithWrongMethod [put]
func MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	// Set the response status to 405 Method Not Allowed
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusMethodNotAllowed)

	// Create a MethodNotAllowedResponse object
	response := ErrorResponse{
		Error:   "Method Not Allowed",
		Message: "The requested method is not allowed for this resource.",
	}

	// Encode the response as JSON and write it to the response writer
	json.NewEncoder(w).Encode(response)
}

func SendErrorResponse(w http.ResponseWriter, errRes ErrorResponse, statusCode int) {
	// Set the response status code
	w.WriteHeader(statusCode)

	// Encode the error response as JSON and write it to the response writer
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(errRes); err != nil {
		http.Error(w, `{"error": "Internal Server Error", "message": "Failed to encode error response."}`, http.StatusInternalServerError)
	}
}

func SendInternalServerError(w http.ResponseWriter, message string) {
	errRes := ErrorResponse{
		Error:   "Internal Server Error",
		Message: message,
	}

	SendErrorResponse(w, errRes, http.StatusInternalServerError)
}
