package model

import (
	"context"

	"firebase.google.com/go/v4/auth"
)

// AuthService defines the interface for authentication services
type AuthService interface {
	// VerifyToken verifies the Firebase ID token and returns the token claims
	VerifyToken(ctx context.Context, idToken string) (*auth.Token, error)

	// GetUser gets a user by their UID
	GetUser(ctx context.Context, uid string) (*auth.UserRecord, error)
}
