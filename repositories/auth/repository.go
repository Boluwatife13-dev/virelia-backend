package auth

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"virelia-backend/models"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrUserExists   = errors.New(
		"email or phone is already registered",
	)
)

type UserRepository struct {
	users *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{users: db.Collection("users")}
}

func (r *UserRepository) Create(
	ctx context.Context, user *models.User,
) error {
	user.CreatedAt = time.Now().UTC()
	user.UpdatedAt = user.CreatedAt

	result, err := r.users.InsertOne(ctx, user)
	if mongo.IsDuplicateKeyError(err) {
		return ErrUserExists
	}
	if err != nil {
		return err
	}

	if id, ok := result.InsertedID.(bson.ObjectID); ok {
		user.ID = id
	}

	return nil
}

func (r *UserRepository) FindByEmail(
	ctx context.Context, email string,
) (*models.User, error) {
	return r.findOne(ctx, bson.M{"email": email})
}

func (r *UserRepository) FindByID(
	ctx context.Context, id bson.ObjectID,
) (*models.User, error) {
	return r.findOne(ctx, bson.M{"_id": id})
}

func (r *UserRepository) RecordSuccessfulLogin(
	ctx context.Context, id bson.ObjectID,
) error {
	now := time.Now().UTC()

	return r.set(ctx, id, bson.M{"last_login_at": now})
}

func (r *UserRepository) MarkVerified(
	ctx context.Context, id bson.ObjectID, field string,
) error {
	return r.set(ctx, id, bson.M{field: true})
}

func (r *UserRepository) UpdatePassword(
	ctx context.Context, id bson.ObjectID, hash string,
) error {
	return r.set(ctx, id, bson.M{"password_hash": hash})
}

func (r *UserRepository) set(
	ctx context.Context, id bson.ObjectID, fields bson.M,
) error {
	fields["updated_at"] = time.Now().UTC()

	_, err := r.users.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": fields},
	)

	return err
}

func (r *UserRepository) findOne(
	ctx context.Context, filter bson.M,
) (*models.User, error) {
	var user models.User

	err := r.users.FindOne(ctx, filter).Decode(&user)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}