package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Blog represents a blog post in the system
type Blog struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title        string             `json:"title" bson:"title"`
	Content      string             `json:"content" bson:"content"`
	AuthorID     string             `json:"author_id" bson:"author_id"`
	Tags         []string           `json:"tags" bson:"tags"`
	ImageURL     string             `json:"image_url" bson:"image_url"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`
	Likes        []string           `json:"likes" bson:"likes"`
	BookmarkedBy []string           `json:"bookmarked_by" bson:"bookmarked_by"`
	IsPublished  bool               `json:"is_published" bson:"is_published"`
}

// NewBlog creates a new blog post
func NewBlog(title, content, authorID string, tags []string, imageURL string) *Blog {
	now := time.Now()
	return &Blog{
		Title:        title,
		Content:      content,
		AuthorID:     authorID,
		Tags:         tags,
		ImageURL:     imageURL,
		CreatedAt:    now,
		UpdatedAt:    now,
		Likes:        []string{},
		BookmarkedBy: []string{},
		IsPublished:  true,
	}
}
