package records

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrNoPatient  = errors.New("patient not found")
	ErrNotFound   = errors.New("record not found")
	ErrNotAllowed = errors.New("you cannot see or change this record")
	ErrOldVersion = errors.New("only the latest version can be corrected")
	ErrTooBig     = errors.New("record details are too long")
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db}
}

const entryColumns = `
	SELECT e.id, e.patient_id, e.facility_id, f.name,
		e.author_id, COALESCE(s.first_name || ' ' || s.last_name, ''),
		e.appointment_id::text, e.kind, e.data, e.version,
		e.replaces_id::text, e.correction_reason, e.created_at
	FROM record_entries e
	JOIN facilities f ON f.id = e.facility_id
	LEFT JOIN staff s ON s.user_id = e.author_id`

const latestOnly = `
	AND NOT EXISTS (
		SELECT 1 FROM record_entries c WHERE c.replaces_id = e.id
	)`

func (r *Repository) list(
	ctx context.Context,
	where string,
	args ...any,
) ([]Entry, error) {
	sql := entryColumns + " WHERE " + where +
		" ORDER BY e.created_at DESC LIMIT 200"

	rows, _ := r.db.Query(ctx, sql, args...)

	return pgx.CollectRows(rows, pgx.RowToStructByPos[Entry])
}

func (r *Repository) get(ctx context.Context, id string) (*Entry, error) {
	list, err := r.list(ctx, "e.id = $1", id)
	if err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return nil, ErrNotFound
	}

	return &list[0], nil
}
