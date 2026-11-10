package keys

import (
	"context"
	"net/http"
)

// A new platform admin has no key yet. Right after
// their password, they add their first key, which
// also logs them in. After that, the key is a must.

func (s *Service) BeginSetup(
	ctx context.Context,
	challenge string,
) (any, error) {
	a, err := s.adminFrom(ctx, challenge, false)
	if err != nil {
		return nil, err
	}

	return s.BeginRegister(ctx, a.id)
}

func (s *Service) FinishSetup(
	ctx context.Context,
	challenge string,
	r *http.Request,
	agent string,
	ip string,
) (any, error) {
	a, err := s.adminFrom(ctx, challenge, false)
	if err != nil {
		return nil, err
	}

	err = s.FinishRegister(ctx, a.id, r)
	if err != nil {
		return nil, err
	}

	return s.tokens.IssueFor(ctx, a.id, agent, ip)
}
