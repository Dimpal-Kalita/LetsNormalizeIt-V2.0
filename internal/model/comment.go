package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Comment represents a comment on a blog post
type Comment struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	BlogID    primitive.ObjectID `json:"blog_id" bson:"blog_id"`
	UserID    string             `json:"user_id" bson:"user_id"`
	Content   string             `json:"content" bson:"content"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
	Likes     []string           `json:"likes" bson:"likes"`
	ParentID  primitive.ObjectID `json:"parent_id,omitempty" bson:"parent_id,omitempty"`
}

// NewComment creates a new comment
func NewComment(blogID primitive.ObjectID, userID, content string, parentID primitive.ObjectID) *Comment {
	now := time.Now()
	return &Comment{
		BlogID:    blogID,
		UserID:    userID,
		Content:   content,
		CreatedAt: now,
		UpdatedAt: now,
		Likes:     []string{},
		ParentID:  parentID,
	}
}
