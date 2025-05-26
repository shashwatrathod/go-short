package dao

import (
	"context" // Added for context propagation
	"database/sql"
	"fmt"
	"time"

	"github.com/shashwatrathod/url-shortner/db"
)

// defines the structure for a short URL record.
type ShortURL struct {
	ShortURL    string    `json:"short_url"`
	OriginalURL string    `json:"original_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// defines the interface for short URL data access operations.
type ShortURLDao interface {
	// creates a new short URL entry in the database.
	CreateShortURL(ctx context.Context, shortUrl string, originalUrl string) (*ShortURL, error)
	
	// retrieves a short URL entry from the database by its short URL.
	FindByShortUrl(ctx context.Context, shortUrl string) (*ShortURL, error)

	// retries a short URL entry from the DB by its original URL.
	FindByOriginalUrl(ctx context.Context, originalUrl string) (*ShortURL, error)
}

// shortURLDaoImpl is the concrete implementation of ShortURLDao.
type shortURLDaoImpl struct {
	connManager *db.ConnectionManager
}

// creates a new instance of ShortURLDaoImpl.
func NewShortURLDao(cm *db.ConnectionManager) ShortURLDao {
	return &shortURLDaoImpl{
		connManager: cm,
	}
}

// creates a new short_url row with provided .
func (d *shortURLDaoImpl) CreateShortURL(ctx context.Context, shortUrl string, originalUrl string) (*ShortURL, error) {
	if d.connManager == nil {
		return nil, fmt.Errorf("ConnectionManager is not initialized in DAO")
	}
	shardDB, err := d.connManager.GetShardByShardKey(shortUrl) // Use shortUrl as sharding key
	if err != nil {
		return nil, fmt.Errorf("failed to get shard for key %s: %w", shortUrl, err)
	}

	query := `INSERT INTO short_urls (short_url, original_url) VALUES ($1, $2)
               RETURNING short_url, original_url, created_at, updated_at`

	var createdURL ShortURL
	err = shardDB.QueryRowContext(ctx, query, shortUrl, originalUrl).Scan(
		&createdURL.ShortURL,
		&createdURL.OriginalURL,
		&createdURL.CreatedAt,
		&createdURL.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create short URL: %w", err)
	}
	return &createdURL, nil
}

// retrieves a short URL entry from the database by its short URL
func (d *shortURLDaoImpl) FindByShortUrl(ctx context.Context, shortUrl string) (*ShortURL, error) {
	if d.connManager == nil {
		return nil, fmt.Errorf("ConnectionManager is not initialized in DAO")
	}
	shardDB, err := d.connManager.GetShardByShardKey(shortUrl) // Use shortUrl as sharding key
	if err != nil {
		return nil, fmt.Errorf("failed to get shard for key %s: %w", shortUrl, err)
	}

	query := `SELECT short_url, original_url, created_at, updated_at FROM short_urls WHERE short_url = $1`

	var fetchedURL ShortURL
	err = shardDB.QueryRowContext(ctx, query, shortUrl).Scan(
		&fetchedURL.ShortURL,
		&fetchedURL.OriginalURL,
		&fetchedURL.CreatedAt,
		&fetchedURL.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("short URL not found: %s", shortUrl) // Consider a custom error type
		}
		return nil, fmt.Errorf("failed to get short URL: %w", err)
	}
	return &fetchedURL, nil
}

// retrieves a short URL entry from the DB by its original URL.
// returns the ShortUrl entry if found, nil otherwise. 
// returns an error if there was an unexpected error in executing the query.
func (d *shortURLDaoImpl) FindByOriginalUrl(ctx context.Context, originalUrl string) (*ShortURL, error) {
	if d.connManager == nil {
		return nil, fmt.Errorf("ConnectionManager is not initialized in DAO")
	}

	// Search across all shards for the original URL
    result, err := d.connManager.ForEachWithResult(func(db *sql.DB) (interface{}, error) {
        query := `SELECT short_url, original_url, created_at, updated_at FROM short_urls WHERE original_url = $1`
        
        var fetchedURL ShortURL
        err := db.QueryRowContext(ctx, query, originalUrl).Scan(
            &fetchedURL.ShortURL,
            &fetchedURL.OriginalURL,
            &fetchedURL.CreatedAt,
            &fetchedURL.UpdatedAt,
        )
        
        if err != nil {
            if err == sql.ErrNoRows {
                return nil, nil // Not found in this shard, continue searching
            }
            return nil, fmt.Errorf("failed to query shard: %w", err)
        }
        
        return fetchedURL, nil
    })

	if err != nil {
        return nil, fmt.Errorf("failed to search for original URL: %w", err)
    }
    
	if result == nil {
		return nil, nil
	}
	
	shortURL := result.(ShortURL)
	return &shortURL, nil
}