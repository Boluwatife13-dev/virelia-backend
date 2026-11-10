package auth

import (
	"context"
	"fmt"
)

func (s *Service) VerifyEmail(ctx context.Context, req VerifyRequest) error {
	user, err := s.repo.FindUser(ctx, "email", normalize(req.Email))
	if err != nil {
		return ErrInvalidCode
	}

	hash := hashCode(req.Code, s.secret)

	err = s.repo.CheckOTP(ctx, user.ID, PurposeVerifyEmail, hash)
	if err != nil {
		return err
	}

	return s.repo.UpdateUser(ctx, user.ID, "email_verified_at = now()")
}

func (s *Service) ResendCode(ctx context.Context, email string) error {
	user, err := s.repo.FindUser(ctx, "email", normalize(email))
	if err != nil || user.EmailVerifiedAt != nil {
		return nil
	}

	return s.sendCode(ctx, user, PurposeVerifyEmail)
}

func (s *Service) sendCode(
	ctx context.Context, user *User, purpose string,
) error {
	code, err := newCode()
	if err != nil {
		return err
	}

	hash := hashCode(code, s.secret)

	if err := s.repo.SaveOTP(ctx, user.ID, purpose, hash); err != nil {
		return err
	}

	body := fmt.Sprintf(
		"<p>Your Virelia code is <b>%s</b>.</p>"+
			"<p>It expires in 10 minutes.</p>", code,
	)

	return s.mailer.Send(ctx, user.Email, "Your Virelia code", body)
}
