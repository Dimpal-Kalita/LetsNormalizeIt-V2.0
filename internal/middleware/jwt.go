package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/dksensei/letsnormalizeit/internal/auth"
	"github.com/dksensei/letsnormalizeit/internal/utils"
	"github.com/gin-gonic/gin"
)

type contextKey string

const UserIDKey contextKey = "userID"
const UserClaimsKey contextKey = "userClaims"

type JWTMiddleware struct {
	authService *auth.Service
}

func NewJWTMiddleware(authService *auth.Service) *JWTMiddleware {
	return &JWTMiddleware{authService: authService}
}

func (m *JWTMiddleware) Authenticate() gin.HandlerFunc {
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

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			logger.Warn("Invalid Authorization header format")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer <token>"})
			return
		}

		tokenString := parts[1]
		claims, err := m.authService.VerifyToken(tokenString)
		if err != nil {
			logger.Warn("Token verification failed: %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		c.Set("uid", claims.UserID)
		c.Set("userEmail", claims.Email)
		c.Set("userName", claims.Name)
		c.Set("isAdmin", claims.IsAdmin)

		ctx := context.WithValue(c.Request.Context(), UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, UserClaimsKey, claims)
		c.Request = c.Request.WithContext(ctx)

		logger.With("userID", claims.UserID).Info("User authenticated successfully")
		c.Next()
	}
}

func (m *JWTMiddleware) RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		isAdmin, exists := c.Get("isAdmin")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			return
		}

		if admin, ok := isAdmin.(bool); !ok || !admin {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			return
		}

		c.Next()
	}
}

func OptionalAuth(authService *auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.Next()
			return
		}

		tokenString := parts[1]
		claims, err := authService.VerifyToken(tokenString)
		if err != nil {
			c.Next()
			return
		}

		c.Set("uid", claims.UserID)
		c.Set("userEmail", claims.Email)
		c.Set("userName", claims.Name)
		c.Set("isAdmin", claims.IsAdmin)

		ctx := context.WithValue(c.Request.Context(), UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, UserClaimsKey, claims)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

func GetUserIDFromContext(ctx context.Context) (string, bool) {
	uid, ok := ctx.Value(UserIDKey).(string)
	return uid, ok
}

func GetClaimsFromContext(ctx context.Context) (*auth.Claims, bool) {
	claims, ok := ctx.Value(UserClaimsKey).(*auth.Claims)
	return claims, ok
}
