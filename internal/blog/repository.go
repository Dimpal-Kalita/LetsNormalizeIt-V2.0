package blog

import (
	"context"

	"github.com/dksensei/letsnormalizeit/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository interface {
	Create(ctx context.Context, blog *model.Blog) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*model.Blog, error)
	List(ctx context.Context, filter bson.M, opts *options.FindOptions) ([]*model.Blog, error)
	Update(ctx context.Context, blog *model.Blog) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	AddLike(ctx context.Context, blogID primitive.ObjectID, userID string) error
	RemoveLike(ctx context.Context, blogID primitive.ObjectID, userID string) error
	AddBookmark(ctx context.Context, blogID primitive.ObjectID, userID string) error
	RemoveBookmark(ctx context.Context, blogID primitive.ObjectID, userID string) error
}

type repository struct {
	db *mongo.Database
}

func NewRepository(db *mongo.Database) Repository {
	return &repository{db: db}
}

func (r *repository) collection() *mongo.Collection {
	return r.db.Collection("blogs")
}

func (r *repository) Create(ctx context.Context, blog *model.Blog) error {
	_, err := r.collection().InsertOne(ctx, blog)
	return err
}

func (r *repository) GetByID(ctx context.Context, id primitive.ObjectID) (*model.Blog, error) {
	var blog model.Blog
	err := r.collection().FindOne(ctx, bson.M{"_id": id}).Decode(&blog)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &blog, err
}

func (r *repository) List(ctx context.Context, filter bson.M, opts *options.FindOptions) ([]*model.Blog, error) {
	cursor, err := r.collection().Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var blogs []*model.Blog
	if err := cursor.All(ctx, &blogs); err != nil {
		return nil, err
	}
	return blogs, nil
}

func (r *repository) Update(ctx context.Context, blog *model.Blog) error {
	_, err := r.collection().ReplaceOne(ctx, bson.M{"_id": blog.ID}, blog)
	return err
}

func (r *repository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection().DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *repository) AddLike(ctx context.Context, blogID primitive.ObjectID, userID string) error {
	_, err := r.collection().UpdateOne(
		ctx,
		bson.M{"_id": blogID},
		bson.M{
			"$addToSet": bson.M{"likes": userID},
		},
	)
	return err
}

func (r *repository) RemoveLike(ctx context.Context, blogID primitive.ObjectID, userID string) error {
	_, err := r.collection().UpdateOne(
		ctx,
		bson.M{"_id": blogID},
		bson.M{
			"$pull": bson.M{"likes": userID},
		},
	)
	return err
}

func (r *repository) AddBookmark(ctx context.Context, blogID primitive.ObjectID, userID string) error {
	_, err := r.collection().UpdateOne(
		ctx,
		bson.M{"_id": blogID},
		bson.M{
			"$addToSet": bson.M{"bookmarked_by": userID},
		},
	)
	return err
}

func (r *repository) RemoveBookmark(ctx context.Context, blogID primitive.ObjectID, userID string) error {
	_, err := r.collection().UpdateOne(
		ctx,
		bson.M{"_id": blogID},
		bson.M{
			"$pull": bson.M{"bookmarked_by": userID},
		},
	)
	return err
}
