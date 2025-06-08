package model

// FirebaseAuthInput represents the input for validating a Firebase token
type FirebaseAuthInput struct {
	Token string `json:"token" binding:"required"`
	User  User   `json:"user" binding:"required"`
}

// AuthResponse represents the response for authentication
type AuthResponse struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}
