package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type OTPPurpose string

const (
	OTPVerifyEmail   OTPPurpose = "verify_email"
	OTPVerifyPhone   OTPPurpose = "verify_phone"
	OTPResetPassword OTPPurpose = "reset_password"
)

type OTP struct {
	ID bson.ObjectID `bson:"_id,omitempty"`

	UserID bson.ObjectID `bson:"user_id"`

	Purpose OTPPurpose `bson:"purpose"`

	CodeHash string `bson:"code_hash"`

	Attempts int `bson:"attempts"`

	ExpiresAt time.Time `bson:"expires_at"`

	CreatedAt time.Time `bson:"created_at"`
}