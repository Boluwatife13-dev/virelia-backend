package audit

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db}
}

const insertEntry = `
	INSERT INTO audit_log (
		actor_id,
		actor_role,
		facility_id,
		action,
		patient_id,
		target_id,
		details,
		ip,
		user_agent
	)
	VALUES (
		NULLIF($1, '')::uuid,
		NULLIF($2, ''),
		NULLIF($3, '')::uuid,
		$4,
		NULLIF($5, '')::uuid,
		NULLIF($6, ''),
		$7,
		NULLIF($8, ''),
		NULLIF($9, '')
	)`

func (r *Repository) Insert(ctx context.Context, e Entry) error {
	if e.Details == nil {
		e.Details = map[string]any{}
	}

	_, err := r.db.Exec(
		ctx,
		insertEntry,
		e.ActorID,
		e.ActorRole,
		e.FacilityID,
		e.Action,
		e.PatientID,
		e.TargetID,
		e.Details,
		e.IP,
		e.UserAgent,
	)

	return err
}
