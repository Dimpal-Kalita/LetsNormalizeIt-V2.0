package cache

import (
	"context"
	"sync"
	"time"
)

type MemoryCache struct {
	mu       sync.RWMutex
	data     map[string]*memoryItem
	enabled  bool
}

type memoryItem struct {
	value      string
	expiration time.Time
}

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		data:    make(map[string]*memoryItem),
		enabled: true,
	}
}

func (c *MemoryCache) Get(_ context.Context, key string) (string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.data[key]
	if !exists {
		return "", ErrCacheMiss
	}

	if time.Now().After(item.expiration) {
		delete(c.data, key)
		return "", ErrCacheMiss
	}

	return item.value, nil
}

func (c *MemoryCache) Set(_ context.Context, key string, value string, exp time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var expiration time.Time
	if exp > 0 {
		expiration = time.Now().Add(exp)
	} else {
		expiration = time.Now().Add(24 * time.Hour)
	}

	c.data[key] = &memoryItem{
		value:      value,
		expiration: expiration,
	}

	return nil
}

func (c *MemoryCache) Delete(_ context.Context, key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.data, key)
	return nil
}

func (c *MemoryCache) DeletePrefix(_ context.Context, prefix string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	for key := range c.data {
		if len(key) >= len(prefix) && key[:len(prefix)] == prefix {
			delete(c.data, key)
		}
	}

	return nil
}

func (c *MemoryCache) Exists(_ context.Context, key string) (bool, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.data[key]
	if !exists {
		return false, nil
	}

	if time.Now().After(item.expiration) {
		delete(c.data, key)
		return false, nil
	}

	return true, nil
}

func (c *MemoryCache) Increment(_ context.Context, key string) (int64, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, exists := c.data[key]
	if !exists {
		c.data[key] = &memoryItem{
			value:      "1",
			expiration: time.Now().Add(24 * time.Hour),
		}
		return 1, nil
	}

	var val int64
	for _, c := range item.value {
		if c >= '0' && c <= '9' {
			val = val*10 + int64(c-'0')
		}
	}

	val++
	item.value = "0"
	for val > 0 {
		item.value = string(rune('0'+val%10)) + item.value
		val /= 10
	}
	if item.value == "0" {
		item.value = "1"
	}

	return int64(len(item.value)), nil
}

func (c *MemoryCache) Close() error {
	return nil
}

func (c *MemoryCache) IsEnabled() bool {
	return c.enabled
}

func (c *MemoryCache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.data)
}

func (c *MemoryCache) Cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for key, item := range c.data {
		if now.After(item.expiration) {
			delete(c.data, key)
		}
	}
}
