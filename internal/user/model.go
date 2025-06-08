package user

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a user in the system
type User struct {
	ID        string               `json:"id" bson:"_id"` // Firebase UID used as MongoDB ID
	Name      string               `json:"name" bson:"name"`
	Email     string               `json:"email" bson:"email"`
	PhotoURL  string               `json:"photo_url" bson:"photo_url"`
	CreatedAt time.Time            `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time            `json:"updated_at" bson:"updated_at"`
	Bookmarks []primitive.ObjectID `json:"bookmarks" bson:"bookmarks"`
	Likes     []primitive.ObjectID `json:"likes" bson:"likes"`
	IsAdmin   bool                 `json:"is_admin" bson:"is_admin"`
}

// NewUser creates a new user from Firebase user information
func NewUser(id, name, email, photoURL string) *User {
	now := time.Now()
	return &User{
		ID:        id,
		Name:      name,
		Email:     email,
		PhotoURL:  photoURL,
		CreatedAt: now,
		UpdatedAt: now,
		Bookmarks: []primitive.ObjectID{},
		Likes:     []primitive.ObjectID{},
		IsAdmin:   false,
	}
}
