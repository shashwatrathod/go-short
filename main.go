package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/shashwatrathod/url-shortner/handlers"
	"github.com/shashwatrathod/url-shortner/middleware"
	"github.com/shashwatrathod/url-shortner/routes"
)


func main() {
	// Initialize router
	router := mux.NewRouter();

	router.Use(middleware.LoggingMiddleware);
	router.Use(middleware.ErrorHandlingMiddleware);

	// Register routes
	routes.RegisterRoutes(router);

	router.NotFoundHandler = http.HandlerFunc(handlers.NotFoundHandler);
	router.MethodNotAllowedHandler = http.HandlerFunc(handlers.MethodNotAllowedHandler);
	
	// Start the server
	http.ListenAndServe(":8080", router);
}