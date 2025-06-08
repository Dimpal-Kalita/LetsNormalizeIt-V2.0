package auth

import (
	"net/http"

	"github.com/dksensei/letsnormalizeit/internal/user"
	"github.com/gin-gonic/gin"
)

// Handler handles HTTP requests related to authentication
type Handler struct {
	authService *Service
	userService *user.Service
}

// NewHandler creates a new auth handler
func NewHandler(authService *Service, userService *user.Service) *Handler {
	return &Handler{
		authService: authService,
		userService: userService,
	}
}

// SignupInput represents the input for signup
type SignupInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required"`
}

// SigninInput represents the input for signin
type SigninInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse represents the response for auth operations
type AuthResponse struct {
	UID   string `json:"uid"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Token string `json:"token,omitempty"`
}

// Register registers the auth routes
func (h *Handler) Register(router *gin.Engine) {
	auth := router.Group("/api/v1/auth")
	{
		auth.POST("/signup", h.Signup)
		auth.POST("/signin", h.Signin)
	}
}

// Signup handles user registration
func (h *Handler) Signup(c *gin.Context) {
	var input SignupInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create user in Firebase and MongoDB
	user, err := h.userService.CreateUser(c.Request.Context(), input.Email, input.Password, input.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// In a real implementation, we would generate a token here
	// using Firebase custom token generation or a similar method
	c.JSON(http.StatusCreated, AuthResponse{
		UID:   user.ID,
		Email: user.Email,
		Name:  user.Name,
	})
}

// Signin handles user login
// Note: In a real-world implementation, the actual login would be handled by Firebase Auth SDK on the client side
// This endpoint is for demonstration purposes only and would not be used in a real app
func (h *Handler) Signin(c *gin.Context) {
	var input SigninInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// In a real implementation, Firebase would handle authentication directly
	// This is just a placeholder to demonstrate the API structure
	c.JSON(http.StatusOK, gin.H{
		"message": "Authentication should be handled by Firebase SDK on the client side",
		"info":    "Call Firebase signInWithEmailAndPassword() from your client app",
	})
}
