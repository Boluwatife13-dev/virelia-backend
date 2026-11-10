package keys

import (
	"context"
	"net/http"
)

// adminFrom reads the step-1 token and makes sure
// it belongs to a platform admin. needKey says if
// they must already have a key (login) or must
// not have one yet (first-time setup).
func (s *Service) adminFrom(
	ctx context.Context,
	challenge string,
	needKey bool,
) (*admin, error) {
	userID, role, err := s.tokens.ParseChallenge(challenge)
	if err != nil {
		return nil, err
	}

	if role != "platform_admin" {
		return nil, ErrNotAdmin
	}

	a, err := s.repo.Admin(ctx, userID)
	if err != nil {
		return nil, err
	}

	hasKey := len(a.keys) > 0

	switch {
	case needKey && !hasKey:
		return nil, ErrNoKey
	case !needKey && hasKey:
		return nil, ErrHasKey
	}

	return a, nil
}

// BeginLogin is step 2 of login (after password).
// The browser gets a challenge for the key to sign.
func (s *Service) BeginLogin(
	ctx context.Context,
	challenge string,
) (any, error) {
	a, err := s.adminFrom(ctx, challenge, true)
	if err != nil {
		return nil, err
	}

	options, session, err := s.wa.BeginLogin(a)
	if err != nil {
		return nil, err
	}

	err = s.repo.SaveChallenge(ctx, a.id, "login", session)

	return options, err
}

// FinishLogin checks the key's signature and, if it
// is good, logs the admin in.
func (s *Service) FinishLogin(
	ctx context.Context,
	challenge string,
	r *http.Request,
	agent string,
	ip string,
) (any, error) {
	a, err := s.adminFrom(ctx, challenge, true)
	if err != nil {
		return nil, err
	}

	session, err := s.repo.TakeChallenge(ctx, a.id, "login")
	if err != nil {
		return nil, err
	}

	cred, err := s.wa.FinishLogin(a, session, r)
	if err != nil || cred.Authenticator.CloneWarning {
		return nil, ErrBadKey
	}

	err = s.repo.UpdateKey(ctx, cred)
	if err != nil {
		return nil, err
	}

	return s.tokens.IssueFor(ctx, a.id, agent, ip)
}
