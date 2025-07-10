package comment

import (
	"context"
	"time"

	"github.com/dksensei/letsnormalizeit/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service interface {
	CreateComment(ctx context.Context, comment *model.Comment) error
	GetCommentByID(ctx context.Context, id primitive.ObjectID) (*model.Comment, error)
	ListCommentsByBlogID(ctx context.Context, blogID primitive.ObjectID, page, limit int64) ([]*model.Comment, error)
	UpdateComment(ctx context.Context, comment *model.Comment) error
	DeleteComment(ctx context.Context, id primitive.ObjectID) error
	LikeComment(ctx context.Context, commentID primitive.ObjectID, userID string) error
	UnlikeComment(ctx context.Context, commentID primitive.ObjectID, userID string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateComment(ctx context.Context, comment *model.Comment) error {
	comment.CreatedAt = time.Now()
	comment.UpdatedAt = comment.CreatedAt
	comment.Likes = []string{}
	return s.repo.Create(ctx, comment)
}

func (s *service) GetCommentByID(ctx context.Context, id primitive.ObjectID) (*model.Comment, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) ListCommentsByBlogID(ctx context.Context, blogID primitive.ObjectID, page, limit int64) ([]*model.Comment, error) {
	opts := options.Find().
		SetSort(bson.M{"created_at": -1}).
		SetSkip((page - 1) * limit).
		SetLimit(limit)
	return s.repo.ListByBlogID(ctx, blogID, opts)
}

func (s *service) UpdateComment(ctx context.Context, comment *model.Comment) error {
	comment.UpdatedAt = time.Now()
	return s.repo.Update(ctx, comment)
}

func (s *service) DeleteComment(ctx context.Context, id primitive.ObjectID) error {
	return s.repo.Delete(ctx, id)
}

func (s *service) LikeComment(ctx context.Context, commentID primitive.ObjectID, userID string) error {
	return s.repo.AddLike(ctx, commentID, userID)
}

func (s *service) UnlikeComment(ctx context.Context, commentID primitive.ObjectID, userID string) error {
	return s.repo.RemoveLike(ctx, commentID, userID)
}
