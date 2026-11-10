package staff

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

const resetTOTP = `
	UPDATE staff s
	SET totp_secret = NULL,
	    totp_enabled_at = NULL
	FROM users u
	WHERE u.id = s.user_id
	AND s.user_id = $1
	AND s.status = 'active'
	AND (
		s.facility_id = NULLIF($2, '')::uuid
		OR ($2 = '' AND s.role = 'hospital_admin')
	)
	RETURNING u.email`

func (r *Repository) ResetTOTP(
	ctx context.Context,
	userID string,
	facilityID string,
) (string, error) {
	var email string

	err := r.db.QueryRow(
		ctx,
		resetTOTP,
		userID,
		facilityID,
	).Scan(&email)

	if errors.Is(err, pgx.ErrNoRows) {
		return "", ErrNotFound
	}

	return email, err
}

func (s *Service) Reset2FA(
	ctx context.Context,
	facilityID string,
	userID string,
	actorID string,
) error {
	if userID == actorID {
		return ErrSelf
	}

	email, err := s.repo.ResetTOTP(ctx, userID, facilityID)
	if err != nil {
		return err
	}

	err = s.sessions.LogoutEverywhere(ctx, userID)
	if err != nil {
		return err
	}

	body := "<p>Your Virelia two-step login was reset.</p>" +
		"<p>Next time you log in, you'll set it up again.</p>" +
		"<p>If you didn't ask for this, tell your admin now.</p>"

	return s.mailer.Send(ctx, email, "Your 2FA was reset", body)
}
