package staff

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"

	"virelia/internal/notify"
)

var ErrSelf = errors.New("you cannot do this to your own account")

type SessionEnder interface {
	LogoutEverywhere(ctx context.Context, userID string) error
}

type Service struct {
	repo     *Repository
	mailer   notify.Brevo
	sessions SessionEnder
	appURL   string
}

func NewService(
	repo *Repository,
	mailer notify.Brevo,
	sessions SessionEnder,
	appURL string,
) *Service {
	return &Service{repo, mailer, sessions, appURL}
}

func (s *Service) List(
	ctx context.Context,
	facilityID string,
) ([]Member, error) {
	return s.repo.List(ctx, facilityID)
}

func (s *Service) Remove(
	ctx context.Context,
	facilityID string,
	userID string,
	actorID string,
) error {
	if userID == actorID {
		return ErrSelf
	}

	err := s.repo.SetStatus(ctx, facilityID, userID, "removed")
	if err != nil {
		return err
	}

	return s.sessions.LogoutEverywhere(ctx, userID)
}

func hash(token string) string {
	sum := sha256.Sum256([]byte(token))

	return hex.EncodeToString(sum[:])
}
