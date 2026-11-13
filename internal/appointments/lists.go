package appointments

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

const listColumns = `
	SELECT
		a.id,
		p.id,
		p.first_name || ' ' || p.last_name,
		p.network_id,
		a.doctor_id,
		s.first_name || ' ' || s.last_name,
		f.id,
		f.name,
		a.starts_at,
		a.ends_at,
		a.status,
		a.reason
	FROM appointments a
	JOIN patients p ON p.id = a.patient_id
	JOIN staff s ON s.user_id = a.doctor_id
	JOIN facilities f ON f.id = a.facility_id`

func collect(rows pgx.Rows) ([]Appointment, error) {
	return pgx.CollectRows(rows, pgx.RowToStructByPos[Appointment])
}


func (r *Repository) forPatientUser(
	ctx context.Context,
	userID string,
) ([]Appointment, error) {
	rows, _ := r.db.Query(ctx, listColumns+`
		WHERE p.user_id = $1
		ORDER BY a.starts_at DESC
		LIMIT 100`, userID)

	return collect(rows)
}


func (r *Repository) onDay(
	ctx context.Context,
	column, value string,
	day time.Time,
) ([]Appointment, error) {
	y, m, d := day.Date()
	from := time.Date(y, m, d, 0, 0, 0, 0, wat)

	rows, _ := r.db.Query(ctx, listColumns+`
		WHERE `+column+` = $1
		AND a.starts_at >= $2
		AND a.starts_at < $3
		ORDER BY a.starts_at`,
		value, from, from.AddDate(0, 0, 1))

	return collect(rows)
}