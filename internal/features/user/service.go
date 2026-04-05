package user

import (
	"context"
	"errors"

	"github.com/dksensei/letsnormalizeit/internal/auth"
	"github.com/dksensei/letsnormalizeit/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Repository interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	UpdateUser(ctx context.Context, id string, update bson.M) error
	DeleteUser(ctx context.Context, id string) error
	AddBookmark(ctx context.Context, userID string, blogID primitive.ObjectID) error
	RemoveBookmark(ctx context.Context, userID string, blogID primitive.ObjectID) error
	AddLike(ctx context.Context, userID string, blogID primitive.ObjectID) error
	RemoveLike(ctx context.Context, userID string, blogID primitive.ObjectID) error
}

type Service struct {
	repo       Repository
	authService auth.AuthService
}

func NewService(repo Repository, authService auth.AuthService) *Service {
	return &Service{
		repo:       repo,
		authService: authService,
	}
}

func (s *Service) Register(ctx context.Context, req *model.RegisterRequest) (*model.AuthResponse, error) {
	return s.authService.Register(ctx, req)
}

func (s *Service) Login(ctx context.Context, req *model.LoginRequest) (*model.AuthResponse, error) {
	return s.authService.Login(ctx, req)
}

func (s *Service) GetUserByID(ctx context.Context, userID string) (*model.User, error) {
	return s.repo.GetUserByID(ctx, userID)
}

func (s *Service) UpdateUserProfile(ctx context.Context, userID, name string) (*model.User, error) {
	err := s.repo.UpdateUser(ctx, userID, bson.M{"name": name})
	if err != nil {
		return nil, err
	}
	return s.repo.GetUserByID(ctx, userID)
}

func (s *Service) ToggleLike(ctx context.Context, userID, blogID string) error {
	blogOID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return errors.New("invalid blog ID")
	}

	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	for _, like := range user.Likes {
		if like == blogOID {
			return s.repo.RemoveLike(ctx, userID, blogOID)
		}
	}

	return s.repo.AddLike(ctx, userID, blogOID)
}

func (s *Service) ToggleBookmark(ctx context.Context, userID, blogID string) error {
	blogOID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return errors.New("invalid blog ID")
	}

	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	for _, bookmark := range user.Bookmarks {
		if bookmark == blogOID {
			return s.repo.RemoveBookmark(ctx, userID, blogOID)
		}
	}

	return s.repo.AddBookmark(ctx, userID, blogOID)
}

func (s *Service) GetBookmarks(ctx context.Context, userID string) ([]primitive.ObjectID, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return user.Bookmarks, nil
}
