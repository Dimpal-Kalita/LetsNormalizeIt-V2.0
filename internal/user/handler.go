package user

import (
	"net/http"

	"github.com/dksensei/letsnormalizeit/internal/model"
	"github.com/dksensei/letsnormalizeit/internal/utils"
	"github.com/gin-gonic/gin"
)

// Handler handles HTTP requests related to users
type Handler struct {
	userService *Service
}

// NewHandler creates a new user handler
func NewHandler(userService *Service) *Handler {
	return &Handler{
		userService: userService,
	}
}

// UserRegistrationInput represents the input for user registration
type UserRegistrationInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	PhotoURL string `json:"photo_url"`
}

// UserResponse represents the response for user operations
type UserResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	PhotoURL  string `json:"photo_url,omitempty"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// Login handles user login after token validation by middleware
func (h *Handler) Login(c *gin.Context) {
	logger := utils.NewLogContext(
		"operation", "Login",
		"path", c.Request.URL.Path,
		"method", c.Request.Method,
		"clientIP", c.ClientIP(),
	)

	logger.Info("Processing user login request")

	// Get the user ID from the context (set by auth middleware)
	uid, exists := c.Get("uid")
	if !exists {
		logger.Error("User ID not found in context - authentication middleware may have failed")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	userID := uid.(string)
	logger.With("userID", userID).Info("User ID extracted from token context")

	// Get user from database
	user, err := h.userService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		logger.With("userID", userID).Warn("User not found in database: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found. Please register first."})
		return
	}

	logger.With("userID", userID, "email", user.Email).Info("User login successful")

	// Return the user data
	c.JSON(http.StatusOK, UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		PhotoURL:  user.PhotoURL,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	})
}

// LoginWithAutoRegister handles user login with automatic registration if user doesn't exist
func (h *Handler) LoginWithAutoRegister(c *gin.Context) {
	logger := utils.NewLogContext(
		"operation", "LoginWithAutoRegister",
		"path", c.Request.URL.Path,
		"method", c.Request.Method,
		"clientIP", c.ClientIP(),
	)

	logger.Info("Processing login with auto-registration request")

	var input UserRegistrationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		logger.Error("Invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.With("email", input.Email, "name", input.Name).Debug("Request body parsed successfully")

	// Get the user ID from the context (set by auth middleware)
	uid, exists := c.Get("uid")
	if !exists {
		logger.Error("User ID not found in context - authentication middleware may have failed")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	userID := uid.(string)
	logger.With("userID", userID).Info("User ID extracted from token context")

	// Try to get existing user first
	user, err := h.userService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		// User doesn't exist, create a new one
		logger.With("userID", userID).Info("User not found in database, creating new user")
		
		newUser := model.NewUser(
			userID,
			input.Name,
			input.Email,
			input.PhotoURL,
		)

		user, err = h.userService.StoreUser(c.Request.Context(), newUser)
		if err != nil {
			logger.With("userID", userID).Error("Failed to create new user: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		logger.With("userID", userID, "email", user.Email).Info("New user created successfully")
	} else {
		logger.With("userID", userID, "email", user.Email).Info("Existing user found, returning user data")
	}

	// Return the user data
	c.JSON(http.StatusOK, UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		PhotoURL:  user.PhotoURL,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	})
}

// RegisterUser handles user registration after token validation by middleware
func (h *Handler) RegisterUser(c *gin.Context) {
	logger := utils.NewLogContext(
		"operation", "RegisterUser",
		"path", c.Request.URL.Path,
		"method", c.Request.Method,
		"clientIP", c.ClientIP(),
	)

	logger.Info("Processing user registration request")

	var input UserRegistrationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		logger.Error("Invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.With("email", input.Email, "name", input.Name).Debug("Request body parsed successfully")

	// Get the user ID from the context (set by auth middleware)
	uid, exists := c.Get("uid")
	if !exists {
		logger.Error("User ID not found in context - authentication middleware may have failed")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	userID := uid.(string)
	logger.With("userID", userID).Info("User ID extracted from token context")

	// Create a new user object
	newUser := model.NewUser(
		userID,
		input.Name,
		input.Email,
		input.PhotoURL,
	)

	// Store the user in MongoDB
	user, err := h.userService.StoreUser(c.Request.Context(), newUser)
	if err != nil {
		logger.With("userID", userID).Error("Failed to store user in database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.With("userID", userID, "email", user.Email).Info("User registration completed successfully")

	// Return the user data
	c.JSON(http.StatusOK, UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		PhotoURL:  user.PhotoURL,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	})
} 