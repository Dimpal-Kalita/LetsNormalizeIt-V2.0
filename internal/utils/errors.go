package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// APIError represents an API error
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ParseFirebaseError parses Firebase error messages
func ParseFirebaseError(err error) (int, string) {
	errMsg := err.Error()

	// Common Firebase auth errors and their corresponding HTTP status codes
	switch {
	case strings.Contains(errMsg, "auth/email-already-exists"):
		return http.StatusConflict, "Email already exists"
	case strings.Contains(errMsg, "auth/invalid-email"):
		return http.StatusBadRequest, "Invalid email format"
	case strings.Contains(errMsg, "auth/user-not-found"):
		return http.StatusNotFound, "User not found"
	case strings.Contains(errMsg, "auth/invalid-credential"):
		return http.StatusUnauthorized, "Invalid credentials"
	case strings.Contains(errMsg, "auth/wrong-password"):
		return http.StatusUnauthorized, "Invalid credentials"
	case strings.Contains(errMsg, "auth/id-token-expired"):
		return http.StatusUnauthorized, "Token expired"
	case strings.Contains(errMsg, "auth/id-token-revoked"):
		return http.StatusUnauthorized, "Token revoked"
	case strings.Contains(errMsg, "auth/invalid-id-token"):
		return http.StatusUnauthorized, "Invalid token"
	case strings.Contains(errMsg, "auth/weak-password"):
		return http.StatusBadRequest, "Password is too weak"
	default:
		log.Printf("Unhandled Firebase error: %v", err)
		return http.StatusInternalServerError, "Authentication error"
	}
}

// FormatJSONError formats an error as JSON
func FormatJSONError(code int, message string) []byte {
	errResponse := APIError{
		Code:    code,
		Message: message,
	}

	jsonBytes, err := json.Marshal(errResponse)
	if err != nil {
		// If JSON marshaling fails, return a simple string
		return []byte(fmt.Sprintf(`{"code":%d,"message":"Error formatting JSON: %s"}`,
			http.StatusInternalServerError, err.Error()))
	}

	return jsonBytes
}

// IsNotFoundError checks if an error is a "not found" error
func IsNotFoundError(err error) bool {
	return err != nil && (strings.Contains(err.Error(), "not found") ||
		strings.Contains(err.Error(), "no documents") ||
		errors.Is(err, errors.New("mongo: no documents in result")))
}
