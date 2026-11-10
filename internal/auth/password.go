package auth

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var ErrWrongPassword = errors.New("current password is incorrect")

func (s *Service) ForgotPassword(ctx context.Context, email string) error {
	user, err := s.repo.FindUser(ctx, "email", normalize(email))
	if err != nil {
		return nil
	}

	return s.sendCode(ctx, user, PurposeResetPassword)
}

func (s *Service) ResetPassword(ctx context.Context, req ResetRequest) error {
	user, err := s.repo.FindUser(ctx, "email", normalize(req.Email))
	if err != nil {
		return ErrInvalidCode
	}

	hash := hashCode(req.Code, s.secret)

	err = s.repo.CheckOTP(ctx, user.ID, PurposeResetPassword, hash)
	if err != nil {
		return err
	}

	return s.setPassword(ctx, user.ID, req.Password)
}

func (s *Service) ChangePassword(
	ctx context.Context, userID string, req ChangePasswordRequest,
) error {
	user, err := s.repo.FindUser(ctx, "id", userID)
	if err != nil {
		return err
	}

	if !passwordMatches(user, req.Current) {
		return ErrWrongPassword
	}

	return s.setPassword(ctx, userID, req.Password)
}

func (s *Service) setPassword(
	ctx context.Context, userID, password string,
) error {
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password), bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	err = s.repo.UpdateUser(
		ctx, userID, "password_hash = $2", string(hash),
	)
	if err != nil {
		return err
	}

	return s.repo.RevokeTokens(ctx, "user_id", userID)
}
