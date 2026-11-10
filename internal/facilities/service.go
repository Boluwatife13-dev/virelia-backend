package facilities

import (
	"context"

	"virelia/internal/notify"
	"virelia/internal/staff"
)

type Inviter interface {
	SendInvite(ctx context.Context, in staff.Invite) error
}

type Service struct {
	repo    *Repository
	mailer  notify.Brevo
	inviter Inviter
}

func NewService(
	repo *Repository,
	mailer notify.Brevo,
	inviter Inviter,
) *Service {
	return &Service{repo, mailer, inviter}
}

func (s *Service) List(
	ctx context.Context,
	status string,
) ([]Facility, error) {
	return s.repo.List(ctx, status)
}
