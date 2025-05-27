package dao

import (
	"context" // Added for context propagation
	"database/sql"
	"fmt"
	"time"

	"github.com/shashwatrathod/url-shortner/db"
)

// defines the structure for a UrlAlias record.
type UrlAlias struct {
	Alias       string    `json:"short_url"`
	OriginalURL string    `json:"original_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// defines the interface for short URL data access operations.
type UrlAliasDao interface {
	// creates a new UrlAlias entry in the database.
	CreateUrlAlias(ctx context.Context, alias string, originalUrl string) (*UrlAlias, error)

	// retrieves a short URL entry from the database by its alias.
	FindByAlias(ctx context.Context, alias string) (*UrlAlias, error)

	// retries a short URL entry from the DB by its original URL.
	FindByOriginalUrl(ctx context.Context, originalUrl string) (*UrlAlias, error)
}

// urlAliasDaoImpl is the concrete implementation of UrlAliasDao.
type urlAliasDaoImpl struct {
	connManager *db.ConnectionManager
}

// creates a new instance of urlAliasDaoImpl.
func NewUrlAliasDao(cm *db.ConnectionManager) UrlAliasDao {
	return &urlAliasDaoImpl{
		connManager: cm,
	}
}

// creates a new url_alias row with provided .
func (d *urlAliasDaoImpl) CreateUrlAlias(ctx context.Context, alias string, originalUrl string) (*UrlAlias, error) {
	if d.connManager == nil {
		return nil, fmt.Errorf("ConnectionManager is not initialized in DAO")
	}
	shardDB, err := d.connManager.GetShardByShardKey(alias) // Use alias as sharding key
	if err != nil {
		return nil, fmt.Errorf("failed to get shard for key %s: %w", alias, err)
	}

	query := `INSERT INTO url_aliases (alias, original_url) VALUES ($1, $2)
               RETURNING alias, original_url, created_at, updated_at`

	var createdUrlAlias UrlAlias
	err = shardDB.QueryRowContext(ctx, query, alias, originalUrl).Scan(
		&createdUrlAlias.Alias,
		&createdUrlAlias.OriginalURL,
		&createdUrlAlias.CreatedAt,
		&createdUrlAlias.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create URL Alias: %w", err)
	}
	return &createdUrlAlias, nil
}

// retrieves a URL Alias entry from the database by its alias
func (d *urlAliasDaoImpl) FindByAlias(ctx context.Context, shortUrl string) (*UrlAlias, error) {
	if d.connManager == nil {
		return nil, fmt.Errorf("ConnectionManager is not initialized in DAO")
	}
	shardDB, err := d.connManager.GetShardByShardKey(shortUrl) // Use alias as sharding key
	if err != nil {
		return nil, fmt.Errorf("failed to get shard for key %s: %w", shortUrl, err)
	}

	query := `SELECT alias, original_url, created_at, updated_at FROM url_aliases WHERE alias = $1`

	var fetchedAlias UrlAlias
	err = shardDB.QueryRowContext(ctx, query, shortUrl).Scan(
		&fetchedAlias.Alias,
		&fetchedAlias.OriginalURL,
		&fetchedAlias.CreatedAt,
		&fetchedAlias.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get Alias: %w", err)
	}
	return &fetchedAlias, nil
}

// retrieves an Alias entry from the DB by its original URL.
// returns the UrlAlias entry if found, nil otherwise.
// returns an error if there was an unexpected error in executing the query.
func (d *urlAliasDaoImpl) FindByOriginalUrl(ctx context.Context, originalUrl string) (*UrlAlias, error) {
	if d.connManager == nil {
		return nil, fmt.Errorf("ConnectionManager is not initialized in DAO")
	}

	// Search across all shards for the original URL
	result, err := d.connManager.ForEachWithResult(func(db *sql.DB) (interface{}, error) {
		query := `SELECT alias, original_url, created_at, updated_at FROM url_aliases WHERE original_url = $1`

		var fetchedAlias UrlAlias
		err := db.QueryRowContext(ctx, query, originalUrl).Scan(
			&fetchedAlias.Alias,
			&fetchedAlias.OriginalURL,
			&fetchedAlias.CreatedAt,
			&fetchedAlias.UpdatedAt,
		)

		if err != nil {
			if err == sql.ErrNoRows {
				return nil, nil // Not found in this shard, continue searching
			}
			return nil, fmt.Errorf("failed to query shard: %w", err)
		}

		return fetchedAlias, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to search for original URL: %w", err)
	}

	if result == nil {
		return nil, nil
	}

	alias := result.(UrlAlias)
	return &alias, nil
}
