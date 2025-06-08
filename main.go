package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"

	_ "github.com/shashwatrathod/url-shortner/docs/swagger"
	"github.com/shashwatrathod/url-shortner/internal/cache"
	"github.com/shashwatrathod/url-shortner/internal/config"
	"github.com/shashwatrathod/url-shortner/internal/db"
	"github.com/shashwatrathod/url-shortner/internal/handlers"
	"github.com/shashwatrathod/url-shortner/internal/middleware"
	"github.com/shashwatrathod/url-shortner/internal/routes"

	httpSwagger "github.com/swaggo/http-swagger"
)

// initializes and returns the db connection manager.
func initDb(conf *config.Config) (*db.ConnectionManager, error) {
	var connConfigs = make([]db.ConnectionConfig, len(conf.DBConfigs))

	for i, dbConfig := range conf.DBConfigs {
		connConfigs[i] = db.ConnectionConfig{
			DSN: fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
				dbConfig.DBUser,
				dbConfig.Password,
				dbConfig.Host,
				dbConfig.Port,
				dbConfig.DBName,
			),
			ShardName: dbConfig.DBName,
		}
	}

	dbManager, err := db.NewConnectionManager(connConfigs)
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize ConnectionManager: %v", err)
	}

	// Apply migrations
	if err := dbManager.ApplyMigrations(); err != nil {
		return nil, fmt.Errorf("Failed to apply migrations: %v", err)
	}

	return dbManager, nil
}

// initializes and returns the redis cache manager
func initRedisCacheManager(ctx context.Context) (cache.CacheManager, error) {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	password := os.Getenv("REDIS_PASSWORD")

	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "6379"
	}

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       0,
	})

	return cache.NewRedisCacheManager(ctx, client)
}

// @title URL Shortener API
// @version 1.0
// @description API Documentation for the Go-Short URL shortening service.

// @host localhost:8080
// @BasePath /api
// @schemes http
func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Error loading .env file, relying on environment variables.")
	}

	conf := config.Load()

	ctx := context.Background()

	// Initialize the DB Connection Manager.
	dbManager, err := initDb(conf)

	if err != nil {
		log.Fatalf("Initializing DBManager : %s", err)
	}

	log.Printf("Initializing DBManager : Success")

	// Initialize Redis Cache Manager
	cacheManager, err := initRedisCacheManager(ctx)
	if err != nil {
		log.Fatalf("Initializing CacheManager : %s", err)
	}

	log.Printf("Initializing CacheManager : Success")

	// Initialize AppEnv
	appEnv := middleware.NewAppEnv(dbManager, cacheManager)

	// Initialize router
	router := mux.NewRouter()

	router.Use(middleware.LoggingMiddleware)
	router.Use(middleware.ErrorHandlingMiddleware)

	router.Use(middleware.ContextMiddleware(appEnv))

	// Register API routes
	routes.RegisterRoutes(router)

	// TODO: Make this route conditional - based on deployment env.
	// Swagger UI route : http://{host}:8080/swagger/index.html
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	router.NotFoundHandler = http.HandlerFunc(handlers.NotFoundHandler)
	router.MethodNotAllowedHandler = http.HandlerFunc(handlers.MethodNotAllowedHandler)

	// Start the server
	log.Println("Starting server on :8080")
	log.Println("Access Swagger at http://localhost:8080/swagger/index.html")
	http.ListenAndServe(":8080", router)
}
