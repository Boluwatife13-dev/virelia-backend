package facilities

import (
	"context"
	"errors"
	"html"
	"strings"
	"time"
)

var ErrBadExpiry = errors.New(
	"licence expiry must be a future date",
)

func (s *Service) Apply(
	ctx context.Context,
	req ApplyRequest,
) (*Facility, error) {
	expiry, err := time.Parse("2006-01-02", req.LicenceExpiresAt)

	if err != nil || expiry.Before(time.Now()) {
		return nil, ErrBadExpiry
	}

	f := &Facility{
		Name:             strings.TrimSpace(req.Name),
		Type:             req.Type,
		LicenceNumber:    strings.TrimSpace(req.LicenceNumber),
		LicenceExpiresAt: expiry,
		State:            req.State,
		Address:          req.Address,
		Email:            strings.ToLower(req.Email),
		Phone:            req.Phone,
		AdminFirstName:   req.AdminFirstName,
		AdminLastName:    req.AdminLastName,
		AdminEmail:       strings.ToLower(req.AdminEmail),
	}

	err = s.repo.Create(ctx, f)
	if err != nil {
		return nil, err
	}

	admins, _ := s.repo.PlatformAdminEmails(ctx)

	name := html.EscapeString(f.Name)
	state := html.EscapeString(f.State)
	body := "<p>" + name + " in " + state +
		" has applied to join Virelia.</p>"

	s.mailer.SendAll(ctx, admins, "New hospital application", body)

	return f, nil
}
