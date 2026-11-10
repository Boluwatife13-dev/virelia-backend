package keys

import (
	"context"
	"errors"

	"github.com/go-webauthn/webauthn/webauthn"
)

var (
	ErrNotAdmin = errors.New("only platform admins use security keys")
	ErrNoKey    = errors.New("no security key registered")
	ErrBadKey   = errors.New("security key check failed")
	ErrHasKey   = errors.New("you already have a key, log in with it")
)

// Tokens comes from the auth module. It reads the
// step-1 login token and issues the real tokens.
type Tokens interface {
	ParseChallenge(token string) (userID, role string, err error)
	IssueFor(ctx context.Context, userID, agent, ip string) (any, error)
}

type Service struct {
	repo   *Repository
	wa     *webauthn.WebAuthn
	tokens Tokens
}

// NewService sets up WebAuthn. rpID is the site's
// domain, e.g. "localhost" or "virelia.ng".
func NewService(
	repo *Repository,
	tokens Tokens,
	rpID string,
	origins []string,
) (*Service, error) {
	wa, err := webauthn.New(&webauthn.Config{
		RPID:          rpID,
		RPDisplayName: "Virelia",
		RPOrigins:     origins,
	})

	return &Service{repo, wa, tokens}, err
}
