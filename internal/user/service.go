package user

import (
	"context"
	"errors"

	"github.com/dksensei/letsnormalizeit/internal/model"
	"github.com/dksensei/letsnormalizeit/internal/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Service handles user-related business logic
type Service struct {
	repo        *Repository
	authService model.AuthService
}

// Ensure Service implements model.UserService
var _ model.UserService = (*Service)(nil)

// NewService creates a new user service
func NewService(repo *Repository, authService model.AuthService) *Service {
	return &Service{
		repo:        repo,
		authService: authService,
	}
}

// GetUserByID gets a user by ID
func (s *Service) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	logger := utils.NewLogContext("userID", id, "operation", "GetUserByID")

	// Try to get from database first
	logger.Debug("Attempting to find user in database")
	user, err := s.repo.FindByID(ctx, id)
	if err == nil {
		logger.Debug("User found in database")
		return user, nil
	}

	// If not found in database, try to get from Firebase
	logger.Debug("User not found in database, trying Firebase")
	firebaseUser, err := s.authService.GetUser(ctx, id)
	if err != nil {
		logger.Error("Failed to get user from Firebase: %v", err)
		return nil, err
	}

	// Create a new user in the database
	logger.Info("Creating new user record from Firebase data")
	newUser := NewUser(
		firebaseUser.UID,
		firebaseUser.DisplayName,
		firebaseUser.Email,
		firebaseUser.PhotoURL,
	)

	if err := s.repo.Create(ctx, newUser); err != nil {
		logger.Error("Failed to create user in database: %v", err)
		return nil, err
	}

	logger.Info("User successfully created in database")
	return newUser, nil
}

// StoreUser stores a user in the database
func (s *Service) StoreUser(ctx context.Context, user *model.User) (*model.User, error) {
	logger := utils.NewLogContext("userID", user.ID, "operation", "StoreUser")

	// Check if user already exists
	existingUser, err := s.repo.FindByID(ctx, user.ID)
	if err == nil {
		// User exists, update the user
		logger.Debug("User exists, updating user")
		existingUser.Name = user.Name
		existingUser.Email = user.Email
		existingUser.PhotoURL = user.PhotoURL

		if err := s.repo.Update(ctx, existingUser); err != nil {
			logger.Error("Failed to update user in database: %v", err)
			return nil, err
		}

		return existingUser, nil
	}

	// User doesn't exist, create a new one
	logger.Debug("User doesn't exist, creating new user")
	if err := s.repo.Create(ctx, user); err != nil {
		logger.Error("Failed to create user in database: %v", err)
		return nil, err
	}

	return user, nil
}

// UpdateUserProfile updates a user's profile
func (s *Service) UpdateUserProfile(ctx context.Context, id, name string) (*model.User, error) {
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

// NewUser creates a new user from Firebase user information
func NewUser(id, name, email, photoURL string) *model.User {
	return model.NewUser(id, name, email, photoURL)
}
