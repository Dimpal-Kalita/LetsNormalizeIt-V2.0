package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/dksensei/letsnormalizeit/internal/utils"
	"github.com/gin-gonic/gin"
)

// RateLimiter implements a simple in-memory rate limiter
type RateLimiter struct {
	// Map of IP/user to last request time and count
	clients map[string]rateLimitClient
	mu      *sync.Mutex
	limit   int           // Maximum number of requests allowed
	window  time.Duration // Time window for rate limiting
}

type rateLimitClient struct {
	count    int       // Number of requests made
	lastSeen time.Time // Last request time
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		clients: make(map[string]rateLimitClient),
		mu:      &sync.Mutex{},
		limit:   limit,
		window:  window,
	}
}

// RateLimit creates a middleware that limits requests based on IP address
func (rl *RateLimiter) RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get client identifier (IP address or user ID if authenticated)
		clientID := c.ClientIP()

		// Create logger with request context
		logger := utils.NewLogContext(
			"path", c.Request.URL.Path,
			"method", c.Request.Method,
			"clientIP", clientID,
		)

		// If user is authenticated, use user ID instead
		if uid, exists := c.Get("uid"); exists {
			clientID = uid.(string)
			logger = logger.With("userID", clientID)
		}

		allow := func() bool {
			rl.mu.Lock()
			defer rl.mu.Unlock()

			now := time.Now()
			client, exists := rl.clients[clientID]

			// If client doesn't exist or window has elapsed, reset the counter
			if !exists || now.Sub(client.lastSeen) > rl.window {
				rl.clients[clientID] = rateLimitClient{
					count:    1,
					lastSeen: now,
				}
				return true
			}

			// If client has exceeded the limit, reject the request
			if client.count >= rl.limit {
				logger.Warn("Rate limit exceeded by client: %s (count: %d, limit: %d)",
					clientID, client.count, rl.limit)
				return false
			}

			// Otherwise, increment the counter and allow the request
			client.count++
			client.lastSeen = now
			rl.clients[clientID] = client
			return true
		}()

		if !allow {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded. Try again later.",
			})
			return
		}

		c.Next()
	}
}

// Cleanup periodically removes old clients from the rate limiter
func (rl *RateLimiter) Cleanup(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			rl.mu.Lock()
			now := time.Now()
			cleanedCount := 0
			for clientID, client := range rl.clients {
				// Remove clients that haven't been seen in a while
				if now.Sub(client.lastSeen) > rl.window*2 {
					delete(rl.clients, clientID)
					cleanedCount++
				}
			}
			totalClients := len(rl.clients)
			rl.mu.Unlock()

			if cleanedCount > 0 {
				utils.Info("Rate limiter cleanup: removed %d stale clients, %d remaining",
					cleanedCount, totalClients)
			}
		}
	}()
}
