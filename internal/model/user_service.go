package model

import (
	"context"
)

// UserService defines the interface for user-related services
type UserService interface {
	// GetUserByID gets a user by ID
	GetUserByID(ctx context.Context, id string) (*User, error)

	// StoreUser stores a user in the database
	StoreUser(ctx context.Context, user *User) (*User, error)

	// UpdateUserProfile updates a user's profile
	UpdateUserProfile(ctx context.Context, id, name string) (*User, error)

	// ToggleBookmark toggles a bookmark for a user
	ToggleBookmark(ctx context.Context, userID, blogID string) error

	// ToggleLike toggles a like for a user
	ToggleLike(ctx context.Context, userID, blogID string) error
}
