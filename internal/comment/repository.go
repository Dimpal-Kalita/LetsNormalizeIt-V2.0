package comment

import (
	"context"

	"github.com/dksensei/letsnormalizeit/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository interface {
	Create(ctx context.Context, comment *model.Comment) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*model.Comment, error)
	ListByBlogID(ctx context.Context, blogID primitive.ObjectID, opts *options.FindOptions) ([]*model.Comment, error)
	Update(ctx context.Context, comment *model.Comment) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	AddLike(ctx context.Context, commentID primitive.ObjectID, userID string) error
	RemoveLike(ctx context.Context, commentID primitive.ObjectID, userID string) error
}

type repository struct {
	db *mongo.Database
}

func NewRepository(db *mongo.Database) Repository {
	return &repository{db: db}
}

func (r *repository) collection() *mongo.Collection {
	return r.db.Collection("comments")
}

func (r *repository) Create(ctx context.Context, comment *model.Comment) error {
	_, err := r.collection().InsertOne(ctx, comment)
	return err
}

func (r *repository) GetByID(ctx context.Context, id primitive.ObjectID) (*model.Comment, error) {
	var comment model.Comment
	err := r.collection().FindOne(ctx, bson.M{"_id": id}).Decode(&comment)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &comment, err
}

func (r *repository) ListByBlogID(ctx context.Context, blogID primitive.ObjectID, opts *options.FindOptions) ([]*model.Comment, error) {
	cursor, err := r.collection().Find(ctx, bson.M{"blog_id": blogID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var comments []*model.Comment
	if err := cursor.All(ctx, &comments); err != nil {
		return nil, err
	}
	return comments, nil
}

func (r *repository) Update(ctx context.Context, comment *model.Comment) error {
	_, err := r.collection().ReplaceOne(ctx, bson.M{"_id": comment.ID}, comment)
	return err
}

func (r *repository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection().DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *repository) AddLike(ctx context.Context, commentID primitive.ObjectID, userID string) error {
	_, err := r.collection().UpdateOne(
		ctx,
		bson.M{"_id": commentID},
		bson.M{
			"$addToSet": bson.M{"likes": userID},
		},
	)
	return err
}

func (r *repository) RemoveLike(ctx context.Context, commentID primitive.ObjectID, userID string) error {
	_, err := r.collection().UpdateOne(
		ctx,
		bson.M{"_id": commentID},
		bson.M{
			"$pull": bson.M{"likes": userID},
		},
	)
	return err
}
