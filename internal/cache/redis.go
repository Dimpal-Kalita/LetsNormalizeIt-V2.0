package cache

import (
	"context"
	"strconv"
	"time"

	"github.com/dksensei/letsnormalizeit/internal/config"
	"github.com/go-redis/redis/v8"
)

type RedisCache struct {
	client  *redis.Client
	enabled bool
}

func NewRedisCache(cfg *config.RedisConfig, enabled bool) (*RedisCache, error) {
	if !enabled {
		return &RedisCache{enabled: false}, nil
	}

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &RedisCache{
		client:  client,
		enabled: true,
	}, nil
}

func (c *RedisCache) Get(ctx context.Context, key string) (string, error) {
	if !c.enabled {
		return "", ErrCacheMiss
	}

	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", ErrCacheMiss
	}
	return val, err
}

func (c *RedisCache) Set(ctx context.Context, key string, value string, exp time.Duration) error {
	if !c.enabled {
		return nil
	}

	return c.client.Set(ctx, key, value, exp).Err()
}

func (c *RedisCache) Delete(ctx context.Context, key string) error {
	if !c.enabled {
		return nil
	}

	return c.client.Del(ctx, key).Err()
}

func (c *RedisCache) DeletePrefix(ctx context.Context, prefix string) error {
	if !c.enabled {
		return nil
	}

	pattern := prefix + "*"
	iter := c.client.Scan(ctx, 0, pattern, 100).Iterator()

	for iter.Next(ctx) {
		if err := c.client.Del(ctx, iter.Val()).Err(); err != nil {
			return err
		}
	}

	return iter.Err()
}

func (c *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	if !c.enabled {
		return false, nil
	}

	result, err := c.client.Exists(ctx, key).Result()
	return result > 0, err
}

func (c *RedisCache) Increment(ctx context.Context, key string) (int64, error) {
	if !c.enabled {
		return 0, nil
	}

	return c.client.Incr(ctx, key).Result()
}

func (c *RedisCache) Close() error {
	if !c.enabled {
		return nil
	}
	return c.client.Close()
}

func (c *RedisCache) IsEnabled() bool {
	return c.enabled
}

func (c *RedisCache) GetInt(ctx context.Context, key string) (int, error) {
	val, err := c.Get(ctx, key)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(val)
}

func (c *RedisCache) SetInt(ctx context.Context, key string, value int, exp time.Duration) error {
	return c.Set(ctx, key, strconv.Itoa(value), exp)
}
