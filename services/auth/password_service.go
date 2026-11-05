package auth

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"

	"virelia-backend/models"
	repo "virelia-backend/repositories/auth"
	"virelia-backend/utils"
)

var ErrWrongPassword = errors.New(
	"current password is incorrect",
)

func (s *Service) ForgotPassword(
	ctx context.Context, email string,
) error {
	user, err := s.users.FindByEmail(ctx, normalizeEmail(email))
	if errors.Is(err, repo.ErrUserNotFound) {
		return nil
	}
	if err != nil {
		return err
	}

	return s.createOTP(ctx, user, models.OTPResetPassword)
}

func (s *Service) ResetPassword(
	ctx context.Context, email, code, password string,
) error {
	user, err := s.users.FindByEmail(ctx, normalizeEmail(email))
	if errors.Is(err, repo.ErrUserNotFound) {
		return ErrInvalidOTP
	}
	if err != nil {
		return err
	}

	err = s.checkOTP(ctx, user.ID, models.OTPResetPassword, code)
	if err != nil {
		return err
	}

	return s.setPassword(ctx, user.ID, password)
}

func (s *Service) ChangePassword(
	ctx context.Context, userID, current, password string,
) error {
	user, err := s.GetUser(ctx, userID)
	if err != nil {
		return err
	}

	if !utils.CheckPassword(current, user.PasswordHash) {
		return ErrWrongPassword
	}

	return s.setPassword(ctx, user.ID, password)
}

func (s *Service) setPassword(
	ctx context.Context, id bson.ObjectID, password string,
) error {
	hash, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	err = s.users.UpdatePassword(ctx, id, hash)
	if err != nil {
		return err
	}

	return s.tokens.RevokeAllForUser(ctx, id)
}