package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// ErrorResponse creates a standardized error response
func ErrorResponse(err error) gin.H {
	if err == nil {
		return gin.H{"error": "Unknown error"}
	}
	return gin.H{"error": err.Error()}
}

// GetQueryInt64 gets an int64 from query parameters with a default value
func GetQueryInt64(c *gin.Context, key string, defaultValue int64) int64 {
	valueStr := c.Query(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.ParseInt(valueStr, 10, 64)
	if err != nil {
		return defaultValue
	}

	return value
}

// GetQueryString gets a string from query parameters with a default value
func GetQueryString(c *gin.Context, key string, defaultValue string) string {
	value := c.Query(key)
	if value == "" {
		return defaultValue
	}
	return value
}
