package auth

import (
	"net/http"

	"github.com/dksensei/letsnormalizeit/internal/model"
	"github.com/gin-gonic/gin"
)

// Handler handles HTTP requests related to authentication
type Handler struct {
	authService *Service
	userService model.UserService
}

// NewHandler creates a new auth handler
func NewHandler(authService *Service, userService model.UserService) *Handler {
	return &Handler{
		authService: authService,
		userService: userService,
	}
}

// FirebaseAuthInput represents the input for validating a Firebase token and creating a user
type FirebaseAuthInput struct {
	Token string     `json:"token" binding:"required"`
	User  model.User `json:"user" binding:"required"`
}

// AuthResponse represents the response for auth operations
type AuthResponse struct {
	UID   string `json:"uid"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

// Register registers the auth routes
func (h *Handler) Register(router *gin.Engine) {
	auth := router.Group("/api/v1/auth")
	{
		auth.POST("/validate-token", h.ValidateToken)
	}
}

// ValidateToken validates a Firebase token and creates/updates the user in MongoDB
func (h *Handler) ValidateToken(c *gin.Context) {
	var input FirebaseAuthInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify the Firebase token
	token, err := h.authService.VerifyToken(c.Request.Context(), input.Token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// Check if the token UID matches the user ID
	if token.UID != input.User.ID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token UID does not match user ID"})
		return
	}

	// Get the existing user or create a new one
	user, err := h.userService.GetUserByID(c.Request.Context(), token.UID)
	if err != nil {
		// User doesn't exist, create a new one
		newUser := model.NewUser(
			token.UID,
			input.User.Name,
			input.User.Email,
			input.User.PhotoURL,
		)

		// Store the user in MongoDB
		user, err = h.userService.StoreUser(c.Request.Context(), newUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, AuthResponse{
		UID:   user.ID,
		Email: user.Email,
		Name:  user.Name,
	})
}
