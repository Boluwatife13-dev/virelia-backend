package staff

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func (r *Repository) CreateInvite(
	ctx context.Context,
	in Invite,
	tokenHash string,
) error {
	_, err := r.db.Exec(ctx, `
        INSERT INTO invites (
          facility_id, email, role, first_name, last_name,
          licence_number, specialty, invited_by,
          token_hash, expires_at
        ) VALUES ($1, $2, $3, $4, $5, NULLIF($6, ''), NULLIF($7, ''),
          NULLIF($8, '')::uuid, $9, now() + interval '72 hours')`,
		in.FacilityID, in.Email, in.Role, in.FirstName,
		in.LastName, in.Licence, in.Specialty, in.InvitedBy,
		tokenHash,
	)

	return err
}

const acceptQuery = `
    WITH i AS (
      UPDATE invites SET accepted_at = now()
      WHERE token_hash = $1
        AND accepted_at IS NULL
        AND expires_at > now()
        AND facility_id IN (
          SELECT id FROM facilities WHERE status = 'approved')
      RETURNING *
    ), u AS (
      INSERT INTO users
        (email, phone, password_hash, email_verified_at)
      SELECT email, $2, $3, now() FROM i
      RETURNING id
    )
    INSERT INTO staff (user_id, facility_id, role, first_name,
      last_name, licence_number, specialty)
    SELECT u.id, i.facility_id, i.role, i.first_name,
      i.last_name, i.licence_number, i.specialty
    FROM u, i
    RETURNING facility_id, role, first_name, last_name`

func (r *Repository) Accept(
	ctx context.Context,
	tokenHash, phone, passwordHash string,
) (*Invite, error) {
	var in Invite

	err := r.db.QueryRow(
		ctx, acceptQuery, tokenHash, phone, passwordHash,
	).Scan(&in.FacilityID, &in.Role, &in.FirstName, &in.LastName)

	var pgErr *pgconn.PgError

	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, ErrBadInvite
	case errors.As(err, &pgErr) && pgErr.Code == "23505":
		return nil, ErrExists
	}

	return &in, err
}
