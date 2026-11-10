package auth

import (
	"context"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func (s *Service) checkLogin(
	ctx context.Context, req LoginRequest,
) (*User, error) {
	user, err := s.repo.FindUser(ctx, "email", normalize(req.Email))
	if err != nil && !errors.Is(err, ErrNotFound) {
		return nil, err
	}

	switch {
	case err != nil || !passwordMatches(user, req.Password):
		return nil, ErrWrongLogin
	case user.Status != "active":
		return nil, ErrInactive
	case user.EmailVerifiedAt == nil:
		return nil, ErrNotVerified
	}

	return user, nil
}

func (s *Service) Login(
	ctx context.Context, req LoginRequest, agent, ip string,
) (*TokenPair, error) {
	user, err := s.checkLogin(ctx, req)
	if err != nil {
		return nil, err
	}

	role, err := s.repo.RoleOf(ctx, user.ID)
	if err != nil || role.Name != RolePatient {
		return nil, ErrWrongLogin
	}

	return s.issueTokens(ctx, user.ID, "", agent, ip)
}

func (s *Service) Refresh(
	ctx context.Context, token, agent, ip string,
) (*TokenPair, error) {
	t, err := s.repo.UseRefreshToken(ctx, hashToken(token))
	if err != nil || time.Now().After(t.ExpiresAt) {
		return nil, ErrInvalidToken
	}

	user, err := s.repo.FindUser(ctx, "id", t.UserID)
	if err != nil || user.Status != "active" {
		return nil, ErrInvalidToken
	}

	return s.issueTokens(ctx, t.UserID, t.FamilyID, agent, ip)
}

func (s *Service) Logout(ctx context.Context, token string) error {
	t, err := s.repo.UseRefreshToken(ctx, hashToken(token))
	if err != nil {
		return nil
	}

	return s.repo.RevokeTokens(ctx, "family_id", t.FamilyID)
}

func passwordMatches(user *User, password string) bool {
	hash, pw := []byte(user.PasswordHash), []byte(password)

	return bcrypt.CompareHashAndPassword(hash, pw) == nil
}
