package cache

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// default expiry - 20 minutes.
const DEFAULT_EXPIRY_SECONDS time.Duration = time.Duration(20) * time.Second

type CacheManager interface {
	// Sets the key with the given value with DEFAULT_EXPIRY_SECONDS.
	// keyStore identifies the bucket where the key is to be stored.
	// for example key=myKey set in the fooStore would not be found in barStore.
	// Overrides the value and resets the time if the key already exists.
	Set(ctx context.Context, keyStore string, key string, value interface{}) error

	// Gets the value for the given key from the given keyStore.
	// Returns nil if the key doesn't exist or is expired.
	Get(ctx context.Context, keyStore string, key string) (interface{}, error)
}

// redisCacheManager is a CacheManager that uses Redis as its cache management engine.
type redisCacheManager struct {
	client *redis.Client
}

func (r *redisCacheManager) Get(ctx context.Context, keyStore string, key string) (interface{}, error) {
	k := fmt.Sprintf("%s:%s", keyStore, key)
	res, err := r.client.Get(ctx, k).Result()

	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	log.Printf("successfully fetched value for the key %s", k)
	return res, nil
}

func (r *redisCacheManager) Set(ctx context.Context, keyStore string, key string, value interface{}) error {
	k := fmt.Sprintf("%s:%s", keyStore, key)
	res, err := r.client.Set(ctx, k, value, DEFAULT_EXPIRY_SECONDS).Result()

	if err != nil {
		return err
	}

	if res != "OK" {
		return fmt.Errorf("error setting the key - %s", k)
	}

	log.Printf("successfully set the key %s", k)
	return nil
}

func NewRedisCacheManager(ctx context.Context, client *redis.Client) (CacheManager, error) {

	if client == nil {
		return nil, fmt.Errorf("Received Nil redis.Client")
	}

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("Unable to Ping Redis. Is the Redis server online?")
	}

	return &redisCacheManager{
		client: client,
	}, nil
}
