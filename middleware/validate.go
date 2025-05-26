package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

// ValidationError defines the structure for validation error responses.
//
// @Description Validation error response structure.
type ValidationError struct {
	Error    string   `json:"error" example:"ValidationError"` // Error type, typically "ValidationError"
	Messages []string `json:"messages"`                        // List of validation error messages
}

func init() {
	validate = validator.New()

	// Register custom tag name function to use json tags
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

func Validate[T any](next func(http.ResponseWriter, *http.Request, T)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload T

		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, `{"error": "Bad Request", "message": "Invalid JSON format"}`, http.StatusBadRequest)
			return
		}

		if err := validate.Struct(payload); err != nil {
			var msg []string
			for _, valErr := range err.(validator.ValidationErrors) {
				fieldName := valErr.Field()
				msg = append(msg, fmt.Sprintf("%s: %s", fieldName, valErr.Tag())) // Simplified message
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&ValidationError{
				Error:    "ValidationError",
				Messages: msg,
			})
			return
		}

		next(w, r, payload)
	}
}
