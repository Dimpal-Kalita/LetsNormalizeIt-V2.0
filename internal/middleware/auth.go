package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/dksensei/letsnormalizeit/internal/auth"
	"github.com/dksensei/letsnormalizeit/internal/utils"
	"github.com/gin-gonic/gin"
)

// Key to store user ID in context
type contextKey string

const UserIDKey contextKey = "userID"

// AuthMiddleware creates a middleware for authenticating requests
func AuthMiddleware(authService *auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := utils.NewLogContext(
			"path", c.Request.URL.Path,
			"method", c.Request.Method,
			"clientIP", c.ClientIP(),
		)

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Warn("Request missing Authorization header")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}
		// Extract the token from the Authorization header
		// Format: Bearer <token>
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			logger.Warn("Invalid Authorization header format")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer <token>"})
			return
		}

		idToken := parts[1]
		token, err := authService.VerifyToken(c.Request.Context(), idToken)
		if err != nil {
			logger.Warn("Token verification failed: %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}
		// Set user ID in the context
		c.Set("uid", token.UID)
		// Add user data to the request context for potential usage in services
		ctx := context.WithValue(c.Request.Context(), UserIDKey, token.UID)
		c.Request = c.Request.WithContext(ctx)

		// Update logger with user information
		logger.With("userID", token.UID).Info("User authenticated successfully")

		c.Next()
	}
}

// GetUserID retrieves the user ID from the context
func GetUserID(ctx context.Context) (string, bool) {
	uid, ok := ctx.Value(UserIDKey).(string)
	return uid, ok
}

// OptionalAuth middleware tries to authenticate but allows requests to proceed if authentication fails
func OptionalAuth(authService *auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next() // Allow the request to proceed without authentication
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Next() // Proceed without setting user ID
			return
		}

		idToken := parts[1]
		token, err := authService.VerifyToken(c.Request.Context(), idToken)
		if err != nil {
			c.Next() // Proceed without setting user ID
			return
		}

		// Set user ID in the context
		c.Set("uid", token.UID)
		ctx := context.WithValue(c.Request.Context(), UserIDKey, token.UID)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

// AdminOnly middleware ensures the user has admin claims
func AdminOnly(authService *auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// First, use the standard auth middleware
		_, exists := c.Get("uid")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			return
		}

		// Check for admin claim in the token
		tokenString := strings.Split(c.GetHeader("Authorization"), " ")[1]
		token, err := authService.VerifyToken(c.Request.Context(), tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// Check if the user has admin claim
		// This assumes that 'admin' is a boolean claim in your Firebase token
		isAdmin, ok := token.Claims["admin"].(bool)
		if !ok || !isAdmin {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			return
		}

		c.Next()
	}
}
