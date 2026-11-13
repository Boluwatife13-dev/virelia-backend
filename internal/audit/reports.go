package audit

import (
	"context"

	"github.com/jackc/pgx/v5"
)

const recordColumns = `
	SELECT
		id,
		actor_id::text,
		actor_role,
		facility_id::text,
		action,
		patient_id::text,
		target_id,
		details,
		ip,
		created_at
	FROM audit_log`

type Filter struct {
	PatientID  string
	ActorID    string
	FacilityID string
}

func (r *Repository) Find(
	ctx context.Context,
	f Filter,
) ([]Record, error) {
	rows, _ := r.db.Query(ctx, recordColumns+`
		WHERE ($1 = '' OR patient_id = NULLIF($1, '')::uuid)
		AND ($2 = '' OR actor_id = NULLIF($2, '')::uuid)
		AND ($3 = '' OR facility_id = NULLIF($3, '')::uuid)
		ORDER BY id DESC
		LIMIT 200`,
		f.PatientID, f.ActorID, f.FacilityID,
	)

	return pgx.CollectRows(rows, pgx.RowToStructByPos[Record])
}

func (r *Repository) ForPatientUser(
	ctx context.Context,
	userID string,
) ([]Record, error) {
	rows, _ := r.db.Query(ctx, recordColumns+`
		WHERE patient_id = (
			SELECT id FROM patients WHERE user_id = $1
		)
		ORDER BY id DESC
		LIMIT 200`,
		userID,
	)

	return pgx.CollectRows(rows, pgx.RowToStructByPos[Record])
}
