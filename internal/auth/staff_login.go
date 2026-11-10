package auth

import (
	"context"
	"time"
)

// StaffLogin is step 1: email and password.
// Platform admins go on to their security key.
// Everyone else goes on to the authenticator app.
func (s *Service) StaffLogin(
	ctx context.Context,
	req LoginRequest,
) (map[string]any, error) {
	user, err := s.checkLogin(ctx, req)
	if err != nil {
		return nil, err
	}

	role, err := s.repo.RoleOf(ctx, user.ID)
	if err != nil || role.Name == RolePatient {
		return nil, ErrWrongLogin
	}

	claims := Claims{Role: role.Name}
	result := map[string]any{}

	if role.Name == RolePlatformAdmin {
		result["use_security_key"] = true
		result["has_key"] = s.repo.HasKey(ctx, user.ID)
	} else {
		err = s.totpStep(user, role, &claims, result)
		if err != nil {
			return nil, err
		}
	}

	token, err := signToken(
		claims,
		user.ID,
		challengeIssuer,
		s.secret,
		5*time.Minute,
	)

	result["challenge_token"] = token

	return result, err
}

// totpStep adds the app setup link on first login.
func (s *Service) totpStep(
	user *User,
	role *Role,
	claims *Claims,
	result map[string]any,
) error {
	result["setup_required"] = !role.TOTPOn

	if role.TOTPOn {
		return nil
	}

	secret, link, err := newTOTP(user.Email)
	if err != nil {
		return err
	}

	claims.TOTPSetup = secret
	result["totp_link"] = link

	return nil
}
