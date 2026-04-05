package cache

import "errors"

var (
	ErrCacheMiss    = errors.New("cache miss")
	ErrCacheFull    = errors.New("cache is full")
	ErrInvalidValue = errors.New("invalid cache value")
)
