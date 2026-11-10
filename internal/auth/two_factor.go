package auth

import (
	"context"
	"errors"
)

var ErrUseKey = errors.New("platform admins must log in with their security key")

type TwoFactorRequest struct {
	ChallengeToken string `json:"challenge_token" binding:"required"`
	Code           string `json:"code" binding:"required,len=6,numeric"`
}

// VerifyTwoFactor is step 2: the 6-digit app code.
func (s *Service) VerifyTwoFactor(
	ctx context.Context,
	req TwoFactorRequest,
	agent string,
	ip string,
) (*TokenPair, error) {
	claims, err := parseToken(
		req.ChallengeToken,
		challengeIssuer,
		s.secret,
	)
	if err != nil {
		return nil, err
	}

	// Platform admins must use a security key.
	if claims.Role == RolePlatformAdmin {
		return nil, ErrUseKey
	}

	secret, err := s.totpSecret(ctx, claims)
	if err != nil || !checkTOTP(secret, req.Code) {
		return nil, ErrInvalidCode
	}

	// First login: save the new secret, encrypted.
	if claims.TOTPSetup != "" {
		sealed, err := s.box.Seal(secret)
		if err != nil {
			return nil, err
		}

		err = s.repo.EnableTOTP(ctx, claims.Subject, sealed)
		if err != nil {
			return nil, err
		}
	}

	return s.issueTokens(ctx, claims.Subject, "", agent, ip)
}

// totpSecret finds the secret to check the code
// against: the new one during setup, or the saved
// one, decrypted.
func (s *Service) totpSecret(
	ctx context.Context,
	claims *Claims,
) (string, error) {
	role, err := s.repo.RoleOf(ctx, claims.Subject)
	if err != nil || role.Name == RolePatient {
		return "", ErrInvalidToken
	}

	if claims.TOTPSetup != "" {
		return claims.TOTPSetup, nil
	}

	if role.TOTPSecret == nil {
		return "", ErrInvalidCode
	}

	return s.box.Open(*role.TOTPSecret)
}
