package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func EnsureIndexes(db *mongo.Database) error {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	users := db.Collection("users")

	_, err := users.Indexes().CreateMany(
		ctx,
		[]mongo.IndexModel{
			{
				Keys: bson.D{{Key: "email", Value: 1}},
				Options: options.Index().
					SetUnique(true).
					SetName("uniq_email"),
			},
			{
				Keys: bson.D{{Key: "phone", Value: 1}},
				Options: options.Index().
					SetUnique(true).
					SetName("uniq_phone"),
			},
		},
	)
	if err != nil {
		return err
	}

	tokens := db.Collection("refresh_tokens")

	_, err = tokens.Indexes().CreateMany(
		ctx,
		[]mongo.IndexModel{
			{
				Keys: bson.D{{Key: "token_hash", Value: 1}},
				Options: options.Index().SetUnique(true),
			},
			{
				Keys: bson.D{{Key: "expires_at", Value: 1}},
				Options: options.Index().SetExpireAfterSeconds(0),
			},
			{Keys: bson.D{{Key: "user_id", Value: 1}}},
			{Keys: bson.D{{Key: "family_id", Value: 1}}},
		},
	)
	if err != nil {
		return err
	}

	otps := db.Collection("otps")

	_, err = otps.Indexes().CreateMany(
		ctx,
		[]mongo.IndexModel{
			{
				Keys: bson.D{
					{Key: "user_id", Value: 1},
					{Key: "purpose", Value: 1},
				},
				Options: options.Index().SetUnique(true),
			},
			{
				Keys: bson.D{{Key: "expires_at", Value: 1}},
				Options: options.Index().SetExpireAfterSeconds(0),
			},
		},
	)

	return err
}