package routes

import (
	"github.com/gorilla/mux"
	"github.com/shashwatrathod/url-shortner/handlers"
)

func RegisterRoutes(router *mux.Router) {
	r := router.PathPrefix("/api").Subrouter();
	r.HandleFunc("/health", handlers.HealthHandler).Methods("GET")
}