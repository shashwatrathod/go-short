package dao

import "time"

type ShortURL struct {
	shortUrl string
	originalUrl string
	createdAt time.Time
	updatedAt time.Time
}

type ShortURLDao interface {
	// CreateShortURL creates a new short URL entry in the database.
	CreateShortURL(shortUrl string, originalUrl string) (ShortURL, error)
	// GetShortURL retrieves a short URL entry from the database by its short URL.
	GetShortURL(shortUrl string) (ShortURL, error)
}

// TODO