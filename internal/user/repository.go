package user

import (
	"context"
	"errors"
	"time"

	"github.com/dksensei/letsnormalizeit/internal/db"
	"github.com/dksensei/letsnormalizeit/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const collectionName = "users"

// Repository handles user data operations
type Repository struct {
	db         *db.MongoDB
	collection string
}

// NewRepository creates a new user repository
func NewRepository(mongodb *db.MongoDB) *Repository {
	return &Repository{
		db:         mongodb,
		collection: collectionName,
	}
}

// FindByID finds a user by ID
func (r *Repository) FindByID(ctx context.Context, id string) (*model.User, error) {
	coll := r.db.GetCollection(r.collection)

	var user model.User
	err := coll.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// FindByEmail finds a user by email
func (r *Repository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	coll := r.db.GetCollection(r.collection)

	var user model.User
	err := coll.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Create creates a new user
func (r *Repository) Create(ctx context.Context, user *model.User) error {
	coll := r.db.GetCollection(r.collection)

	// Check if user already exists
	count, err := coll.CountDocuments(ctx, bson.M{"_id": user.ID})
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("user already exists")
	}

	_, err = coll.InsertOne(ctx, user)
	return err
}

// Update updates an existing user
func (r *Repository) Update(ctx context.Context, user *model.User) error {
	coll := r.db.GetCollection(r.collection)

	user.UpdatedAt = time.Now()

	_, err := coll.ReplaceOne(ctx, bson.M{"_id": user.ID}, user)
	return err
}

// AddBookmark adds a bookmark to a user
func (r *Repository) AddBookmark(ctx context.Context, userID string, blogID primitive.ObjectID) error {
	coll := r.db.GetCollection(r.collection)

	_, err := coll.UpdateOne(
		ctx,
		bson.M{"_id": userID},
		bson.M{
			"$addToSet": bson.M{"bookmarks": blogID},
			"$set":      bson.M{"updated_at": time.Now()},
		},
	)
	return err
}

// RemoveBookmark removes a bookmark from a user
func (r *Repository) RemoveBookmark(ctx context.Context, userID string, blogID primitive.ObjectID) error {
	coll := r.db.GetCollection(r.collection)

	_, err := coll.UpdateOne(
		ctx,
		bson.M{"_id": userID},
		bson.M{
			"$pull": bson.M{"bookmarks": blogID},
			"$set":  bson.M{"updated_at": time.Now()},
		},
	)
	return err
}

// AddLike adds a like to a user
func (r *Repository) AddLike(ctx context.Context, userID string, blogID primitive.ObjectID) error {
	coll := r.db.GetCollection(r.collection)

	_, err := coll.UpdateOne(
		ctx,
		bson.M{"_id": userID},
		bson.M{
			"$addToSet": bson.M{"likes": blogID},
			"$set":      bson.M{"updated_at": time.Now()},
		},
	)
	return err
}

// RemoveLike removes a like from a user
func (r *Repository) RemoveLike(ctx context.Context, userID string, blogID primitive.ObjectID) error {
	coll := r.db.GetCollection(r.collection)

	_, err := coll.UpdateOne(
		ctx,
		bson.M{"_id": userID},
		bson.M{
			"$pull": bson.M{"likes": blogID},
			"$set":  bson.M{"updated_at": time.Now()},
		},
	)
	return err
}
