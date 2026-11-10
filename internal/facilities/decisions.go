package facilities

import (
	"context"
	"html"

	"virelia/internal/staff"
)

func (s *Service) Check(
	ctx context.Context,
	id string,
	adminID string,
	note string,
) error {
	_, err := s.repo.Check(ctx, id, adminID, note)

	return err
}

func (s *Service) Approve(
	ctx context.Context,
	id string,
	adminID string,
) error {
	f, err := s.repo.Approve(ctx, id, adminID)
	if err != nil {
		return err
	}

	invite := staff.Invite{
		FacilityID: f.ID,
		Email:      f.AdminEmail,
		Role:       "hospital_admin",
		FirstName:  f.AdminFirstName,
		LastName:   f.AdminLastName,
		InvitedBy:  adminID,
	}

	return s.inviter.SendInvite(ctx, invite)
}

func (s *Service) Reject(
	ctx context.Context,
	id string,
	adminID string,
	reason string,
) error {
	f, err := s.repo.Reject(ctx, id, adminID, reason)
	if err != nil {
		return err
	}

	name := html.EscapeString(f.Name)
	why := html.EscapeString(reason)

	body := "<p>Your application for " + name +
		" was not approved.</p>" +
		"<p>Reason: " + why + "</p>"

	return s.mailer.Send(ctx, f.AdminEmail, "Your Virelia application", body)
}
