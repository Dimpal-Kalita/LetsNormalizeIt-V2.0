package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           string               `json:"id" bson:"_id"`
	Name         string               `json:"name" bson:"name"`
	Email        string               `json:"email" bson:"email"`
	PasswordHash string               `json:"-" bson:"password_hash"`
	PhotoURL     string               `json:"photo_url" bson:"photo_url"`
	CreatedAt    time.Time            `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time            `json:"updated_at" bson:"updated_at"`
	Bookmarks    []primitive.ObjectID `json:"bookmarks" bson:"bookmarks"`
	Likes        []primitive.ObjectID `json:"likes" bson:"likes"`
	IsAdmin      bool                 `json:"is_admin" bson:"is_admin"`
}

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

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type RegisterRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=100"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type AuthResponse struct {
	User      *User      `json:"user"`
	TokenPair *TokenPair `json:"tokens"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
