package auth

import (
	"context"
	"errors"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"

	"github.com/dksensei/letsnormalizeit/internal/config"
	"github.com/dksensei/letsnormalizeit/internal/model"
	"google.golang.org/api/option"
)

// Service handles Firebase authentication
type Service struct {
	client *auth.Client
}

// Ensure Service implements model.AuthService
var _ model.AuthService = (*Service)(nil)

// NewService creates a new instance of the Firebase auth service
func NewService(config *config.FirebaseConfig) (*Service, error) {
	opt := option.WithCredentialsFile(config.CredentialsFile)
	app, err := firebase.NewApp(context.Background(), &firebase.Config{
		ProjectID: config.ProjectID,
	}, opt)
	if err != nil {
		return nil, err
	}

	client, err := app.Auth(context.Background())
	if err != nil {
		return nil, err
	}

	return &Service{
		client: client,
	}, nil
}

// VerifyToken verifies the Firebase ID token and returns the token claims
func (s *Service) VerifyToken(ctx context.Context, idToken string) (*auth.Token, error) {
	if idToken == "" {
		return nil, errors.New("id token is empty")
	}

	token, err := s.client.VerifyIDToken(ctx, idToken)
	if err != nil {
		log.Printf("Error verifying ID token: %v\n", err)
		return nil, err
	}

	return token, nil
}

// GetUser gets a user by their UID
func (s *Service) GetUser(ctx context.Context, uid string) (*auth.UserRecord, error) {
	return s.client.GetUser(ctx, uid)
}
