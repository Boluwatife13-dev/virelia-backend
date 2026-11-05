package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Role string

const (
	RolePatient       Role = "patient"
	RoleDoctor        Role = "doctor"
	RoleNurse         Role = "nurse"
	RoleHospitalAdmin Role = "hospital_admin"
	RolePlatformAdmin Role = "platform_admin"
)

type AccountStatus string

const (
	StatusPending     AccountStatus = "pending"
	StatusActive      AccountStatus = "active"
	StatusSuspended   AccountStatus = "suspended"
	StatusDeactivated AccountStatus = "deactivated"
)

type User struct {
	ID            bson.ObjectID `bson:"_id,omitempty" json:"id"`
	FirstName     string        `bson:"first_name" json:"first_name"`
	LastName      string        `bson:"last_name" json:"last_name"`
	Email         string        `bson:"email" json:"email"`
	Phone         string        `bson:"phone" json:"phone"`
	PasswordHash  string        `bson:"password_hash" json:"-"`
	Role          Role          `bson:"role" json:"role"`
	Status        AccountStatus `bson:"status" json:"status"`
	EmailVerified bool          `bson:"email_verified" json:"email_verified"`
	PhoneVerified bool          `bson:"phone_verified" json:"phone_verified"`

	FailedLoginAttempts int        `bson:"failed_login_attempts" json:"-"`
	LockedUntil         *time.Time `bson:"locked_until,omitempty" json:"-"`
	LastLoginAt         *time.Time `bson:"last_login_at,omitempty" json:"last_login_at,omitempty"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}