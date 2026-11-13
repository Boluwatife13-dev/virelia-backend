package appointments

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type block struct {
	Starts int64
	Ends   int64
}

func (r *Repository) hours(
	ctx context.Context,
	doctorID string,
	weekday int,
) ([]block, error) {
	rows, _ := r.db.Query(ctx, `
		SELECT
			extract(epoch FROM starts)::bigint,
			extract(epoch FROM ends)::bigint
		FROM doctor_hours
		WHERE doctor_id = $1 AND weekday = $2`,
		doctorID, weekday)

	return pgx.CollectRows(rows, pgx.RowToStructByPos[block])
}

func (r *Repository) busy(
	ctx context.Context,
	doctorID string,
	from, to time.Time,
) ([]Slot, error) {
	rows, _ := r.db.Query(ctx, `
		SELECT starts_at, ends_at
		FROM appointments
		WHERE doctor_id = $1
		AND status <> 'cancelled'
		AND starts_at < $3
		AND ends_at > $2`,
		doctorID, from, to)

	return pgx.CollectRows(rows, pgx.RowToStructByPos[Slot])
}
