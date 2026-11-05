package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type RefreshToken struct {
	ID bson.ObjectID `bson:"_id,omitempty"`

	UserID bson.ObjectID `bson:"user_id"`

	FamilyID string `bson:"family_id"`

	TokenHash string `bson:"token_hash"`

	UserAgent string `bson:"user_agent"`

	IPAddress string `bson:"ip_address"`

	ExpiresAt time.Time `bson:"expires_at"`

	RevokedAt *time.Time `bson:"revoked_at,omitempty"`

	CreatedAt time.Time `bson:"created_at"`
}