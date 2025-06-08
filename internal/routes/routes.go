package routes

import (
	"github.com/gorilla/mux"
	"github.com/shashwatrathod/url-shortner/internal/handlers"
	"github.com/shashwatrathod/url-shortner/internal/middleware"
)

func RegisterRoutes(router *mux.Router) {
	r := router.PathPrefix("/api").Subrouter()
	r.HandleFunc("/health", handlers.HealthHandler).Methods("GET")
	r.HandleFunc("/create", middleware.Validate(handlers.CreateUrlAliasHandler)).Methods("POST")
	r.HandleFunc("/{alias}", handlers.GetUrlAliasHandler).Methods("GET")
}
