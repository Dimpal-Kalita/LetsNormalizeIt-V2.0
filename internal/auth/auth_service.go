package auth

import (
	"context"
	"time"

	"github.com/dksensei/letsnormalizeit/internal/model"
)

type AuthService interface {
	Register(ctx context.Context, req *model.RegisterRequest) (*model.AuthResponse, error)
	Login(ctx context.Context, req *model.LoginRequest) (*model.AuthResponse, error)
	RefreshTokens(ctx context.Context, refreshToken string) (*model.TokenPair, error)
	GetUserByID(ctx context.Context, userID string) (*model.User, error)
}

type authService struct {
	jwtService  *Service
	userRepo    Repository
}

func NewAuthService(jwtService *Service, userRepo Repository) AuthService {
	return &authService{
		jwtService:  jwtService,
		userRepo:    userRepo,
	}
}

func (s *authService) Register(ctx context.Context, req *model.RegisterRequest) (*model.AuthResponse, error) {
	existing, _ := s.userRepo.GetUserByEmail(ctx, req.Email)
	if existing != nil {
		return nil, ErrEmailExists
	}

	passwordHash, err := s.jwtService.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		ID:           GenerateUserID(),
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: passwordHash,
		CreatedAt:    timeNow(),
		UpdatedAt:    timeNow(),
	}

	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	tokenPair, err := s.jwtService.GenerateTokenPair(user.ID, user.Email, user.Name, user.IsAdmin)
	if err != nil {
		return nil, err
	}

	return &model.AuthResponse{
		User: user,
		TokenPair: &model.TokenPair{
			AccessToken:  tokenPair.AccessToken,
			RefreshToken: tokenPair.RefreshToken,
			ExpiresIn:    tokenPair.ExpiresIn,
			TokenType:    "Bearer",
		},
	}, nil
}

func (s *authService) Login(ctx context.Context, req *model.LoginRequest) (*model.AuthResponse, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, ErrInvalidEmail
	}

	if !s.jwtService.CheckPassword(req.Password, user.PasswordHash) {
		return nil, ErrInvalidPassword
	}

	tokenPair, err := s.jwtService.GenerateTokenPair(user.ID, user.Email, user.Name, user.IsAdmin)
	if err != nil {
		return nil, err
	}

	return &model.AuthResponse{
		User: user,
		TokenPair: &model.TokenPair{
			AccessToken:  tokenPair.AccessToken,
			RefreshToken: tokenPair.RefreshToken,
			ExpiresIn:    tokenPair.ExpiresIn,
			TokenType:    "Bearer",
		},
	}, nil
}

func (s *authService) RefreshTokens(ctx context.Context, refreshToken string) (*model.TokenPair, error) {
	tokenPair, err := s.jwtService.RefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	return &model.TokenPair{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
		TokenType:    "Bearer",
	}, nil
}

func (s *authService) GetUserByID(ctx context.Context, userID string) (*model.User, error) {
	return s.userRepo.GetUserByID(ctx, userID)
}

func timeNow() time.Time {
	return time.Now()
}
