package auth

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

const (
	RolePatient       = "patient"
	RoleDoctor        = "doctor"
	RoleFrontDesk     = "front_desk"
	RoleHospitalAdmin = "hospital_admin"
	RolePlatformAdmin = "platform_admin"
)

type Role struct {
	Name       string
	FacilityID string
	TOTPSecret *string
	TOTPOn     bool
}

func (r *Repository) RoleOf(
	ctx context.Context, userID string,
) (*Role, error) {
	var role Role
	var facility *string

	err := r.db.QueryRow(ctx, `
		SELECT s.role, s.facility_id::text,
		       s.totp_secret, s.totp_enabled_at IS NOT NULL
		FROM staff s
		LEFT JOIN facilities f ON f.id = s.facility_id
		WHERE s.user_id = $1
		  AND s.status = 'active'
		  AND (s.facility_id IS NULL OR f.status = 'approved')`,
		userID,
	).Scan(&role.Name, &facility, &role.TOTPSecret, &role.TOTPOn)

	if errors.Is(err, pgx.ErrNoRows) {
		return r.patientRole(ctx, userID)
	}
	if err != nil {
		return nil, err
	}

	if facility != nil {
		role.FacilityID = *facility
	}

	return &role, nil
}

func (r *Repository) patientRole(
	ctx context.Context, userID string,
) (*Role, error) {
	var exists bool

	err := r.db.QueryRow(ctx,
		"SELECT EXISTS (SELECT 1 FROM patients WHERE user_id = $1)",
		userID,
	).Scan(&exists)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, ErrNotFound
	}

	return &Role{Name: RolePatient}, nil
}

func (r *Repository) EnableTOTP(
	ctx context.Context, userID, secret string,
) error {
	_, err := r.db.Exec(ctx, `
		UPDATE staff
		SET totp_secret = $2, totp_enabled_at = now()
		WHERE user_id = $1`, userID, secret)

	return err
}

func (s *Service) LogoutEverywhere(
	ctx context.Context, userID string,
) error {
	return s.repo.RevokeTokens(ctx, "user_id", userID)
}
