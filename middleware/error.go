package middleware

import (
	"log"
	"net/http"
	"runtime/debug"
)

// ErrorHandlingMiddleware recovers from panics, logs the error, and sends a 500 response.
func ErrorHandlingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                // Log the panic and stack trace
                log.Printf("PANIC: %v\n%s", err, debug.Stack())

                // Respond with a 500 Internal Server Error
                // Avoid sending detailed error information to the client in production
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            }
        }()

        // Call the next handler in the chain
        next.ServeHTTP(w, r)
    })
}