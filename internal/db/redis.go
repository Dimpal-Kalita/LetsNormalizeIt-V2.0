package db

import (
	"context"
	"log"

	"github.com/dksensei/letsnormalizeit/internal/config"
	"github.com/go-redis/redis/v8"
)

// Redis represents a Redis connection
type Redis struct {
	Client *redis.Client
}

// NewRedis creates a new Redis connection
func NewRedis(cfg *config.RedisConfig) (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// Ping Redis to verify the connection
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	log.Println("Connected to Redis")
	return &Redis{
		Client: client,
	}, nil
}

// Close closes the Redis connection
func (r *Redis) Close() error {
	return r.Client.Close()
}
