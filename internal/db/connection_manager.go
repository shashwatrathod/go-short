package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/shashwatrathod/url-shortner/internal/utils"
)

type ConnectionManager struct {
	shards       []*sql.DB
	shardsByName map[string]int
}

type ConnectionConfig struct {
	DSN       string
	ShardName string
}

// initializes a new ConnectionManager by opening connections configured in the provided configs.
func NewConnectionManager(configs []ConnectionConfig) (*ConnectionManager, error) {
	shards := make([]*sql.DB, len(configs))
	shardsByName := make(map[string]int)

	for idx, config := range configs {
		db, err := sql.Open("postgres", config.DSN)
		if err != nil {
			return nil, err
		}

		err = db.Ping()
		if err != nil {
			return nil, err
		}

		shards[idx] = db
		shardsByName[config.ShardName] = idx
		log.Printf("Connected to shard: %s", config.ShardName)
	}
	return &ConnectionManager{shards: shards, shardsByName: shardsByName}, nil
}

// returns the DB Shard responsible to handle the provided key.
func (cm *ConnectionManager) GetShardByShardKey(key string) (*sql.DB, error) {

	if key == "" {
		return nil, fmt.Errorf("The key is empty.")
	}

	hash := utils.Hash(key)
	shardIdx := hash % uint64(len(cm.shards))

	return cm.shards[shardIdx], nil
}

// iterates over all database connections and executes the provided function
// Returns the first non-nil result, or nil if no results found
func (cm *ConnectionManager) ForEachWithResult(fn func(db *sql.DB) (interface{}, error)) (interface{}, error) {
	for _, db := range cm.shards {
		result, err := fn(db)
		if err != nil {
			return nil, err
		}
		if result != nil {
			return result, nil
		}
	}
	return nil, nil
}

// applies db migrations from the DB_MIGRATION_DIR to all the db shards
// using goose.
// returns error if any of the migrations couldn't be applied.
func (cm *ConnectionManager) ApplyMigrations() error {

	migrationsDir := os.Getenv("DB_MIGRATION_DIR")
	if migrationsDir == "" {
		return fmt.Errorf("DB_MIGRATION_DIR environment variable not set")
	}

	dbDriver := os.Getenv("DB_DRIVER")
	if dbDriver == "" {
		return fmt.Errorf("DB_DRIVER environment variable not set")
	}

	log.Printf("applying DB Migrations on all shards using Goose..")

	// configure Goose.
	if err := goose.SetDialect(os.Getenv("DB_DRIVER")); err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}
	log.Printf("successfully set goose dialect.")

	for shardName, dbIdx := range cm.shardsByName {
		db := cm.shards[dbIdx]

		if err := goose.Up(db, migrationsDir); err != nil {
			return fmt.Errorf("failed to apply migrations to shard %s: %w", shardName, err)
		}
		log.Printf("successfully applied migrations to shard: %s", shardName)
	}

	log.Println("all shards processed for migrations.")
	return nil
}

// closes all connections held by this connection manager.
func (cm *ConnectionManager) CloseAll() {
	for _, db := range cm.shards {
		if err := db.Close(); err != nil {
			log.Printf("error closing database connection: %v", err)
		}
	}
}
