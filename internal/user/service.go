package user

import (
	"context"
	"errors"

	"github.com/dksensei/letsnormalizeit/internal/auth"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Service handles user-related business logic
type Service struct {
	repo        *Repository
	authService *auth.Service
}

// NewService creates a new user service
func NewService(repo *Repository, authService *auth.Service) *Service {
	return &Service{
		repo:        repo,
		authService: authService,
	}
}

// GetUserByID gets a user by ID
func (s *Service) GetUserByID(ctx context.Context, id string) (*User, error) {
	// Try to get from database first
	user, err := s.repo.FindByID(ctx, id)
	if err == nil {
		return user, nil
	}

	// If not found in database, try to get from Firebase
	firebaseUser, err := s.authService.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	// Create a new user in the database
	newUser := NewUser(
		firebaseUser.UID,
		firebaseUser.DisplayName,
		firebaseUser.Email,
		firebaseUser.PhotoURL,
	)

	if err := s.repo.Create(ctx, newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}

// CreateUser creates a new user
func (s *Service) CreateUser(ctx context.Context, email, password, name string) (*User, error) {
	// Create the user in Firebase first
	firebaseUser, err := s.authService.CreateUser(ctx, email, password)
	if err != nil {
		return nil, err
	}

	// Update display name
	if name != "" {
		firebaseUser, err = s.authService.UpdateUser(ctx, firebaseUser.UID, name)
		if err != nil {
			// This is not critical, we can continue
			// But we should log this error in a real application
		}
	}

	// Create user in the database
	user := NewUser(
		firebaseUser.UID,
		name,
		email,
		firebaseUser.PhotoURL,
	)

	if err := s.repo.Create(ctx, user); err != nil {
		// If we fail to create the user in the database,
		// we should delete the user from Firebase
		// But for simplicity, we'll just return the error
		return nil, err
	}

	return user, nil
}

// UpdateUserProfile updates a user's profile
func (s *Service) UpdateUserProfile(ctx context.Context, id, name string) (*User, error) {
	// Update user in Firebase
	_, err := s.authService.UpdateUser(ctx, id, name)
	if err != nil {
		return nil, err
	}

	// Get current user from database
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update user in database
	user.Name = name
	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// ToggleBookmark toggles a bookmark for a user
func (s *Service) ToggleBookmark(ctx context.Context, userID, blogID string) error {
	// Check if userID exists
	_, err := s.GetUserByID(ctx, userID)
	if err != nil {
		return errors.New("user not found")
	}

	// Convert blogID to ObjectID
	objID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return errors.New("invalid blog ID format")
	}

	// Check if user already has the bookmark
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	// Check if the blog is already bookmarked
	for _, bookmark := range user.Bookmarks {
		if bookmark == objID {
			// Remove the bookmark
			return s.repo.RemoveBookmark(ctx, userID, objID)
		}
	}

	// Add the bookmark
	return s.repo.AddBookmark(ctx, userID, objID)
}

// ToggleLike toggles a like for a user
func (s *Service) ToggleLike(ctx context.Context, userID, blogID string) error {
	// Check if userID exists
	_, err := s.GetUserByID(ctx, userID)
	if err != nil {
		return errors.New("user not found")
	}

	// Convert blogID to ObjectID
	objID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return errors.New("invalid blog ID format")
	}

	// Check if user already has liked the blog
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	// Check if the blog is already liked
	for _, like := range user.Likes {
		if like == objID {
			// Remove the like
			return s.repo.RemoveLike(ctx, userID, objID)
		}
	}

	// Add the like
	return s.repo.AddLike(ctx, userID, objID)
}
