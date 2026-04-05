package cache

import (
	"context"
	"testing"
	"time"
)

func TestMemoryCache(t *testing.T) {
	cache := NewMemoryCache()
	ctx := context.Background()

	t.Run("Set and Get", func(t *testing.T) {
		err := cache.Set(ctx, "key1", "value1", 1*time.Hour)
		if err != nil {
			t.Fatalf("Set() error = %v", err)
		}

		val, err := cache.Get(ctx, "key1")
		if err != nil {
			t.Fatalf("Get() error = %v", err)
		}

		if val != "value1" {
			t.Errorf("Get() = %s, want %s", val, "value1")
		}
	})

	t.Run("Get non-existent key", func(t *testing.T) {
		_, err := cache.Get(ctx, "nonexistent")
		if err != ErrCacheMiss {
			t.Errorf("Get() error = %v, want %v", err, ErrCacheMiss)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		cache.Set(ctx, "key2", "value2", 1*time.Hour)
		cache.Delete(ctx, "key2")

		_, err := cache.Get(ctx, "key2")
		if err != ErrCacheMiss {
			t.Errorf("Get() after Delete() error = %v, want %v", err, ErrCacheMiss)
		}
	})

	t.Run("Exists", func(t *testing.T) {
		cache.Set(ctx, "key3", "value3", 1*time.Hour)

		exists, err := cache.Exists(ctx, "key3")
		if err != nil {
			t.Fatalf("Exists() error = %v", err)
		}
		if !exists {
			t.Error("Exists() = false, want true")
		}

		exists, err = cache.Exists(ctx, "nonexistent")
		if err != nil {
			t.Fatalf("Exists() error = %v", err)
		}
		if exists {
			t.Error("Exists() = true, want false")
		}
	})

	t.Run("DeletePrefix", func(t *testing.T) {
		cache.Set(ctx, "user:1", "data1", 1*time.Hour)
		cache.Set(ctx, "user:2", "data2", 1*time.Hour)
		cache.Set(ctx, "post:1", "data3", 1*time.Hour)

		cache.DeletePrefix(ctx, "user:")

		_, err := cache.Get(ctx, "user:1")
		if err != ErrCacheMiss {
			t.Errorf("Get(user:1) after DeletePrefix() error = %v, want %v", err, ErrCacheMiss)
		}

		_, err = cache.Get(ctx, "user:2")
		if err != ErrCacheMiss {
			t.Errorf("Get(user:2) after DeletePrefix() error = %v, want %v", err, ErrCacheMiss)
		}

		val, _ := cache.Get(ctx, "post:1")
		if val != "data3" {
			t.Errorf("Get(post:1) = %s, want %s", val, "data3")
		}
	})

	t.Run("Expiration", func(t *testing.T) {
		cache.Set(ctx, "temp", "value", 50*time.Millisecond)

		_, err := cache.Get(ctx, "temp")
		if err != nil {
			t.Errorf("Get() immediately after Set() error = %v", err)
		}

		time.Sleep(100 * time.Millisecond)

		_, err = cache.Get(ctx, "temp")
		if err != ErrCacheMiss {
			t.Errorf("Get() after expiration error = %v, want %v", err, ErrCacheMiss)
		}
	})

	t.Run("IsEnabled", func(t *testing.T) {
		if !cache.IsEnabled() {
			t.Error("IsEnabled() = false, want true")
		}
	})

	t.Run("Close", func(t *testing.T) {
		if err := cache.Close(); err != nil {
			t.Errorf("Close() error = %v", err)
		}
	})
}

func TestMemoryCacheCleanup(t *testing.T) {
	cache := NewMemoryCache()
	ctx := context.Background()

	cache.Set(ctx, "key1", "value1", 50*time.Millisecond)
	cache.Set(ctx, "key2", "value2", 1*time.Hour)

	time.Sleep(100 * time.Millisecond)

	cache.Cleanup()

	_, err := cache.Get(ctx, "key1")
	if err != ErrCacheMiss {
		t.Errorf("Get(key1) after Cleanup() error = %v, want %v", err, ErrCacheMiss)
	}

	_, err = cache.Get(ctx, "key2")
	if err != nil {
		t.Errorf("Get(key2) after Cleanup() error = %v", err)
	}
}
