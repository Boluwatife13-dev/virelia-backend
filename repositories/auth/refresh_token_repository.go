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
	ErrTokenNotFound = errors.New("refresh token not found")
	ErrTokenRevoked  = errors.New("refresh token already revoked")
)

type RefreshTokenRepository struct {
	tokens *mongo.Collection
}

func NewRefreshTokenRepository(
	db *mongo.Database,
) *RefreshTokenRepository {
	return &RefreshTokenRepository{
		tokens: db.Collection("refresh_tokens"),
	}
}

func (r *RefreshTokenRepository) Create(
	ctx context.Context, token *models.RefreshToken,
) error {
	token.CreatedAt = time.Now().UTC()

	_, err := r.tokens.InsertOne(ctx, token)

	return err
}

func (r *RefreshTokenRepository) FindByHash(
	ctx context.Context, tokenHash string,
) (*models.RefreshToken, error) {
	var token models.RefreshToken

	filter := bson.M{"token_hash": tokenHash}

	err := r.tokens.FindOne(ctx, filter).Decode(&token)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrTokenNotFound
	}
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (r *RefreshTokenRepository) Revoke(
	ctx context.Context, id bson.ObjectID,
) error {
	result, err := r.tokens.UpdateOne(
		ctx,
		bson.M{"_id": id, "revoked_at": nil},
		revokeUpdate(),
	)
	if err == nil && result.ModifiedCount == 0 {
		return ErrTokenRevoked
	}

	return err
}

func (r *RefreshTokenRepository) RevokeFamily(
	ctx context.Context, familyID string,
) error {
	return r.revokeMany(ctx, bson.M{"family_id": familyID})
}

func (r *RefreshTokenRepository) RevokeAllForUser(
	ctx context.Context, userID bson.ObjectID,
) error {
	return r.revokeMany(ctx, bson.M{"user_id": userID})
}

func (r *RefreshTokenRepository) revokeMany(
	ctx context.Context, filter bson.M,
) error {
	filter["revoked_at"] = nil

	_, err := r.tokens.UpdateMany(ctx, filter, revokeUpdate())

	return err
}

func revokeUpdate() bson.M {
	return bson.M{
		"$set": bson.M{"revoked_at": time.Now().UTC()},
	}
}