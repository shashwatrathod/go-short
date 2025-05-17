package middleware

import (
	"log"
	"net/http"
	"time"
)

// responseWriter is a wrapper around http.ResponseWriter to capture the status code
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.status = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

// LoggingMiddleware logs the incoming HTTP request and its response status.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a responseWriter to capture the status code
		rw := &responseWriter{ResponseWriter: w, status: http.StatusOK} // Default to 200 OK

		// Call the next handler in the chain
		next.ServeHTTP(rw, r)

		duration := time.Since(start)

		// Log the request details and response status
		log.Printf(
			"[%s] %s %s %d %dms",
			r.Method,
			r.RequestURI,
			r.Proto,
			rw.status,
			duration.Milliseconds(),
		)
	})
}
