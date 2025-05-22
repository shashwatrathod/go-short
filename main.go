package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/shashwatrathod/url-shortner/db"
	"github.com/shashwatrathod/url-shortner/handlers"
	"github.com/shashwatrathod/url-shortner/middleware"
	"github.com/shashwatrathod/url-shortner/routes"
)


func main() {
	// Load .env file
    err := godotenv.Load()
    if err != nil {
        log.Println("Warning: Error loading .env file, relying on environment variables.")
    }

	// Initialize ConnectionManager
	configs := []db.ConnectionConfig{
		{
			DSN: fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			"urls",
			), 
			ShardName: "urls",
		},
		{
			DSN: fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			"urls_2",
			), 
			ShardName: "urls_1",
		},
		{
			DSN: fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			"urls_2",
			), 
			ShardName: "urls_2",
		},
	}

	dbManager, err := db.NewConnectionManager(configs)
	if err != nil {
        log.Fatalf("Failed to initialize ConnectionManager: %v", err)
    }
	defer dbManager.CloseAll()

	// Apply migrations
    if err := dbManager.ApplyMigrations(); err != nil {
        log.Fatalf("Failed to apply migrations: %v", err)
    }
    log.Println("Database migrations applied successfully.")

	// Initialize router
	router := mux.NewRouter();

	router.Use(middleware.LoggingMiddleware);
	router.Use(middleware.ErrorHandlingMiddleware);

	appEnv := middleware.NewAppEnv(dbManager)
	router.Use(middleware.ContextMiddleware(appEnv))

	// Register routes
	routes.RegisterRoutes(router);

	router.NotFoundHandler = http.HandlerFunc(handlers.NotFoundHandler);
	router.MethodNotAllowedHandler = http.HandlerFunc(handlers.MethodNotAllowedHandler);
	
	// Start the server
	log.Println("Starting server on :8080")
	http.ListenAndServe(":8080", router);
}