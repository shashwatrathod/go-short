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

type ValidationError struct {
    Error   string `json:"error"`
    Messages []string `json:"messages"`
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
		var payload T;

        if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
            http.Error(w, "Invalid JSON format", http.StatusBadRequest)
            return
        }
		
		if err := validate.Struct(payload); err != nil {
			var msg []string;
			for  _, valErr := range err.(validator.ValidationErrors) {
				msg = append(msg, fmt.Sprintf("%s: %s", valErr.Field(), valErr.Error()));
			}

			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&ValidationError{
				Error: "ValidationError",
				Messages: msg,
			})
			return
		}

		next(w, r, payload)
	}
}