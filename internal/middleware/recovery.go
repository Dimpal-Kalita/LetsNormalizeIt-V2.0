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
				// Log the error and stack trace
				utils.Error("Panic recovered: %v\nStack trace: %s", err, debug.Stack())

				// Return an error response
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": fmt.Sprintf("Internal server error: %v", err),
				})
			}
		}()

		c.Next()
	}
}
