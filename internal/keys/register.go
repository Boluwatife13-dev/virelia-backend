package keys

import (
	"context"
	"net/http"
)

// BeginRegister is step 1 of adding a key. The
// browser gets options to pass to the key.
func (s *Service) BeginRegister(
	ctx context.Context,
	userID string,
) (any, error) {
	a, err := s.repo.Admin(ctx, userID)
	if err != nil {
		return nil, err
	}

	options, session, err := s.wa.BeginRegistration(a)
	if err != nil {
		return nil, err
	}

	err = s.repo.SaveChallenge(ctx, userID, "register", session)

	return options, err
}

// FinishRegister is step 2. The browser sends back
// what the key signed, and we save the key.
func (s *Service) FinishRegister(
	ctx context.Context,
	userID string,
	r *http.Request,
) error {
	session, err := s.repo.TakeChallenge(ctx, userID, "register")
	if err != nil {
		return err
	}

	a, err := s.repo.Admin(ctx, userID)
	if err != nil {
		return err
	}

	cred, err := s.wa.FinishRegistration(a, session, r)
	if err != nil {
		return ErrBadKey
	}

	return s.repo.SaveKey(ctx, userID, cred)
}
