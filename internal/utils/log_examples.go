package utils

import (
	"time"
)

// LogExample demonstrates various structured logging patterns
func LogExample() {
	// Basic logging
	Info("This is a basic info message")
	Error("This is an error message: %v", "something went wrong")

	// Logging with fields
	With("userID", "u123", "requestID", "req456").Info("Processing user request")

	// Creating a logger with context
	logger := NewLogContext("module", "authentication", "serverID", "srv1")
	logger.Info("Starting authentication module")

	// Adding more context to an existing logger
	userLogger := logger.With("userID", "u789")
	userLogger.Debug("User login attempt")

	// Logging complex data structures
	userLogger.Info("User data processed",
		"metadata", map[string]interface{}{
			"loginTime": time.Now(),
			"ipAddress": "192.168.1.1",
			"userAgent": "Mozilla/5.0...",
		},
	)

	// Error logging with context
	if err := someOperationThatFails(); err != nil {
		userLogger.Error("Operation failed: %v", err)
	}

	// You can also use this pattern in your middlewares, handlers, and services
	// to keep track of request context through the entire request lifecycle.
}

// someOperationThatFails is just for the example
func someOperationThatFails() error {
	return nil // simulated success
}
