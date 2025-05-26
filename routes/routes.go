package routes

import (
	"github.com/gorilla/mux"
	"github.com/shashwatrathod/url-shortner/handlers"
	"github.com/shashwatrathod/url-shortner/middleware"
)

func RegisterRoutes(router *mux.Router) {
	r := router.PathPrefix("/api").Subrouter();
	r.HandleFunc("/health", handlers.HealthHandler).Methods("GET")
	r.HandleFunc("/create", middleware.Validate(handlers.CreateShortUrlHandler)).Methods("POST")
	r.HandleFunc("/{shortUrl}", handlers.GetShortUrlHandler).Methods("GET")
}