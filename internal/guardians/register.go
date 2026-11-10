package guardians

import (
	"context"
	"strings"
)

func (s *Service) RegisterByFrontDesk(
	ctx context.Context,
	req FrontDeskChild,
	staffID string,
	facilityID string,
) (string, error) {
	c, err := s.newChild(req.Child)
	if err != nil {
		return "", err
	}

	c.AddedBy = staffID
	c.FacilityID = facilityID
	c.Verified = true

	guardian := strings.ToUpper(strings.TrimSpace(req.GuardianNetworkID))

	err = s.repo.AddChild(ctx, "network_id", guardian, c)

	return c.NetworkID, err
}

func (s *Service) AddMyChild(
	ctx context.Context,
	req Child,
	userID string,
) (string, error) {
	c, err := s.newChild(req)
	if err != nil {
		return "", err
	}

	c.AddedBy = userID

	err = s.repo.AddChild(ctx, "user_id", userID, c)

	return c.NetworkID, err
}

func (s *Service) MyChildren(
	ctx context.Context,
	userID string,
) ([]Link, error) {
	return s.repo.MyChildren(ctx, userID)
}

func (s *Service) GuardiansOf(
	ctx context.Context,
	childNetworkID string,
) ([]Link, error) {
	id := strings.ToUpper(strings.TrimSpace(childNetworkID))

	return s.repo.GuardiansOf(ctx, id)
}

func (s *Service) Verify(ctx context.Context, linkID string) error {
	return s.repo.Verify(ctx, linkID)
}

func (s *Service) Revoke(ctx context.Context, linkID string) error {
	return s.repo.Revoke(ctx, linkID)
}

func (s *Service) Leave(ctx context.Context, linkID, userID string) error {
	return s.repo.Leave(ctx, linkID, userID)
}
