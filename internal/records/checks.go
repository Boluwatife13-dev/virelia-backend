package records

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

func (r *Repository) treats(
	ctx context.Context,
	doctorID, facilityID, patientID string,
) (bool, error) {
	var ok bool

	err := r.db.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM appointments
			WHERE doctor_id = $1
			AND facility_id = $2
			AND patient_id = $3
			AND status <> 'cancelled'
		)`, doctorID, facilityID, patientID,
	).Scan(&ok)

	return ok, err
}

func (r *Repository) isLatest(ctx context.Context, id string) (bool, error) {
	var corrected bool

	err := r.db.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM record_entries WHERE replaces_id = $1
		)`, id,
	).Scan(&corrected)

	return !corrected, err
}

func (r *Repository) patientByNetworkID(
	ctx context.Context,
	networkID string,
) (string, error) {
	var id string

	err := r.db.QueryRow(ctx,
		`SELECT id FROM patients WHERE network_id = $1`,
		networkID,
	).Scan(&id)

	if errors.Is(err, pgx.ErrNoRows) {
		return "", ErrNoPatient
	}

	return id, err
}
