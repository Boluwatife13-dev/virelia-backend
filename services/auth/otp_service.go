package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"

	"virelia-backend/models"
	repo "virelia-backend/repositories/auth"
	"virelia-backend/utils"
)

var (
	ErrInvalidOTP     = errors.New("invalid or expired code")
	ErrSMSUnavailable = errors.New("phone codes are not available yet")
)

var verifiedFields = map[models.OTPPurpose]string{
	models.OTPVerifyEmail: "email_verified",
	models.OTPVerifyPhone: "phone_verified",
}

func (s *Service) SendVerification(
	ctx context.Context, userID string, purpose models.OTPPurpose,
) error {
	if purpose == models.OTPVerifyPhone {
		return ErrSMSUnavailable
	}

	user, err := s.GetUser(ctx, userID)
	if err != nil {
		return err
	}

	return s.createOTP(ctx, user, purpose)
}

func (s *Service) VerifyContact(
	ctx context.Context, userID string,
	purpose models.OTPPurpose, code string,
) error {
	id, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	if err := s.checkOTP(ctx, id, purpose, code); err != nil {
		return err
	}

	return s.users.MarkVerified(ctx, id, verifiedFields[purpose])
}

func (s *Service) createOTP(
	ctx context.Context, user *models.User, purpose models.OTPPurpose,
) error {
	code, err := utils.GenerateOTP()
	if err != nil {
		return err
	}

	err = s.otps.Save(ctx, &models.OTP{
		UserID:    user.ID,
		Purpose:   purpose,
		CodeHash:  utils.HashOTP(code, s.jwtSecret),
		ExpiresAt: time.Now().UTC().Add(utils.OTPTTL),
	})
	if err != nil {
		return err
	}

	html := fmt.Sprintf(
		"<p>Your Virelia code is <b>%s</b>.</p>"+
			"<p>It expires in 10 minutes.</p>",
		code,
	)

	return s.mailer.Send(ctx, user.Email, "Your Virelia code", html)
}

func (s *Service) checkOTP(
	ctx context.Context, userID bson.ObjectID,
	purpose models.OTPPurpose, code string,
) error {
	otp, err := s.otps.Find(ctx, userID, purpose)
	if err == nil && time.Now().After(otp.ExpiresAt) {
		return ErrInvalidOTP
	}
	if err == nil {
		err = s.otps.UseAttempt(ctx, otp.ID, utils.MaxOTPAttempts)
	}
	if errors.Is(err, repo.ErrOTPNotFound) {
		return ErrInvalidOTP
	}
	if err != nil {
		return err
	}

	if !utils.CheckOTP(code, otp.CodeHash, s.jwtSecret) {
		return ErrInvalidOTP
	}

	return s.otps.Delete(ctx, otp.ID)
}