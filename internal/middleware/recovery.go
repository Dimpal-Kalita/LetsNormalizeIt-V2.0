package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/dksensei/letsnormalizeit/internal/utils"
	"github.com/gin-gonic/gin"
)

// Recovery returns a middleware that recovers from any panics
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Create a logger with request context
				logger := utils.With(
					"path", c.Request.URL.Path,
					"method", c.Request.Method,
					"clientIP", c.ClientIP(),
					"userAgent", c.Request.UserAgent(),
				)

				// Log the error and stack trace with context
				logger.Error("Panic recovered: %v\nStack trace: %s", err, debug.Stack())

				// Return an error response
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": fmt.Sprintf("Internal server error: %v", err),
				})
			}
		}()

		c.Next()
	}
}
