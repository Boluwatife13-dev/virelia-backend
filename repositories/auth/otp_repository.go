package auth

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"virelia-backend/models"
)

var ErrOTPNotFound = errors.New("otp not found")

type OTPRepository struct {
	otps *mongo.Collection
}

func NewOTPRepository(db *mongo.Database) *OTPRepository {
	return &OTPRepository{otps: db.Collection("otps")}
}

func (r *OTPRepository) Save(
	ctx context.Context, otp *models.OTP,
) error {
	otp.CreatedAt = time.Now().UTC()

	filter := bson.M{
		"user_id": otp.UserID,
		"purpose": otp.Purpose,
	}

	_, err := r.otps.ReplaceOne(
		ctx,
		filter,
		otp,
		options.Replace().SetUpsert(true),
	)

	return err
}

func (r *OTPRepository) Find(
	ctx context.Context,
	userID bson.ObjectID,
	purpose models.OTPPurpose,
) (*models.OTP, error) {
	filter := bson.M{
		"user_id": userID,
		"purpose": purpose,
	}

	var otp models.OTP

	err := r.otps.FindOne(ctx, filter).Decode(&otp)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrOTPNotFound
	}

	return &otp, err
}


func (r *OTPRepository) UseAttempt(
	ctx context.Context,
	id bson.ObjectID,
	maxAttempts int,
) error {
	filter := bson.M{
		"_id":      id,
		"attempts": bson.M{"$lt": maxAttempts},
	}

	result, err := r.otps.UpdateOne(
		ctx,
		filter,
		bson.M{"$inc": bson.M{"attempts": 1}},
	)
	if err != nil {
		return err
	}
	if result.ModifiedCount == 0 {
		return ErrOTPNotFound
	}

	return nil
}

func (r *OTPRepository) Delete(
	ctx context.Context, id bson.ObjectID,
) error {
	_, err := r.otps.DeleteOne(ctx, bson.M{"_id": id})

	return err
}