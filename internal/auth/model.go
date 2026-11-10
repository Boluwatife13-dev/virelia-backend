package auth

import "time"

const (
	PurposeVerifyEmail   = "verify_email"
	PurposeResetPassword = "reset_password"
)

type User struct {
	ID              string     `json:"id"`
	Email           string     `json:"email"`
	Phone           string     `json:"phone"`
	PasswordHash    string     `json:"-"`
	Status          string     `json:"status"`
	EmailVerifiedAt *time.Time `json:"email_verified_at"`
	CreatedAt       time.Time  `json:"created_at"`
}

type Patient struct {
	ID             string    `json:"id"`
	UserID         string    `json:"user_id"`
	NetworkID      string    `json:"network_id"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	DateOfBirth    time.Time `json:"date_of_birth"`
	Sex            string    `json:"sex"`
	State          *string   `json:"state"`
	Address        *string   `json:"address"`
	NINHash        *string   `json:"-"`
	EmergencyName  *string   `json:"-"`
	EmergencyPhone *string   `json:"-"`
}

type OTP struct {
	ID        string
	UserID    string
	CodeHash  string
	Attempts  int
	ExpiresAt time.Time
}

type RefreshToken struct {
	ID        string
	UserID    string
	FamilyID  string
	ExpiresAt time.Time
	RevokedAt *time.Time
}

type RegisterRequest struct {
	FirstName      string `json:"first_name" binding:"required,max=100"`
	LastName       string `json:"last_name" binding:"required,max=100"`
	Email          string `json:"email" binding:"required,email"`
	Phone          string `json:"phone" binding:"required,e164"`
	Password       string `json:"password" binding:"required,min=8,max=72"`
	DateOfBirth    string `json:"date_of_birth" binding:"required,datetime=2006-01-02"`
	Sex            string `json:"sex" binding:"required,oneof=male female"`
	State          string `json:"state" binding:"omitempty,max=50"`
	Address        string `json:"address" binding:"omitempty,max=255"`
	NIN            string `json:"nin" binding:"required,len=11,numeric"`
	EmergencyName  string `json:"emergency_contact_name" binding:"omitempty,max=100"`
	EmergencyPhone string `json:"emergency_contact_phone" binding:"omitempty,e164"`

	AcceptTerms bool `json:"accept_terms" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type EmailRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type VerifyRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required,len=6,numeric"`
}

type ResetRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Code     string `json:"code" binding:"required,len=6,numeric"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}

type ChangePasswordRequest struct {
	Current  string `json:"current_password" binding:"required"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}

type TokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
