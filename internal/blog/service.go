package blog

import (
	"context"
	"time"

	"github.com/dksensei/letsnormalizeit/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service interface {
	CreateBlog(ctx context.Context, blog *model.Blog) error
	GetBlogByID(ctx context.Context, id primitive.ObjectID) (*model.Blog, error)
	ListBlogs(ctx context.Context, filter bson.M, page, limit int64) ([]*model.Blog, error)
	UpdateBlog(ctx context.Context, blog *model.Blog) error
	DeleteBlog(ctx context.Context, id primitive.ObjectID) error
	LikeBlog(ctx context.Context, blogID primitive.ObjectID, userID string) error
	UnlikeBlog(ctx context.Context, blogID primitive.ObjectID, userID string) error
	BookmarkBlog(ctx context.Context, blogID primitive.ObjectID, userID string) error
	UnbookmarkBlog(ctx context.Context, blogID primitive.ObjectID, userID string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateBlog(ctx context.Context, blog *model.Blog) error {
	blog.CreatedAt = time.Now()
	blog.UpdatedAt = blog.CreatedAt
	blog.Likes = []string{}
	blog.BookmarkedBy = []string{}
	return s.repo.Create(ctx, blog)
}

func (s *service) GetBlogByID(ctx context.Context, id primitive.ObjectID) (*model.Blog, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) ListBlogs(ctx context.Context, filter bson.M, page, limit int64) ([]*model.Blog, error) {
	opts := options.Find().
		SetSort(bson.M{"created_at": -1}).
		SetSkip((page - 1) * limit).
		SetLimit(limit)
	return s.repo.List(ctx, filter, opts)
}

func (s *service) UpdateBlog(ctx context.Context, blog *model.Blog) error {
	blog.UpdatedAt = time.Now()
	return s.repo.Update(ctx, blog)
}

func (s *service) DeleteBlog(ctx context.Context, id primitive.ObjectID) error {
	return s.repo.Delete(ctx, id)
}

func (s *service) LikeBlog(ctx context.Context, blogID primitive.ObjectID, userID string) error {
	return s.repo.AddLike(ctx, blogID, userID)
}

func (s *service) UnlikeBlog(ctx context.Context, blogID primitive.ObjectID, userID string) error {
	return s.repo.RemoveLike(ctx, blogID, userID)
}

func (s *service) BookmarkBlog(ctx context.Context, blogID primitive.ObjectID, userID string) error {
	return s.repo.AddBookmark(ctx, blogID, userID)
}

func (s *service) UnbookmarkBlog(ctx context.Context, blogID primitive.ObjectID, userID string) error {
	return s.repo.RemoveBookmark(ctx, blogID, userID)
}
