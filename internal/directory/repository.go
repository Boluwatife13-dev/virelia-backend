package directory

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db}
}

const searchDoctors = `
	SELECT
		s.user_id,
		s.first_name,
		s.last_name,
		s.specialty,
		f.id,
		f.name,
		f.state,
		p.fee_naira,
		p.slot_minutes,
		p.bio
	FROM staff s
	JOIN doctor_profiles p ON p.user_id = s.user_id
	JOIN facilities f ON f.id = s.facility_id
	WHERE s.role = 'doctor'
	AND s.status = 'active'
	AND f.status = 'approved'
	AND p.accepting
	AND ($1 = '' OR s.specialty ILIKE '%' || $1 || '%')
	AND ($2 = '' OR f.state ILIKE $2)
	AND ($3 = '' OR f.id = NULLIF($3, '')::uuid)
	AND ($4 = '' OR s.first_name || ' ' || s.last_name ILIKE '%' || $4 || '%')
	ORDER BY s.last_name, s.first_name
	LIMIT 50`

func (r *Repository) Search(
	ctx context.Context,
	q Search,
) ([]Doctor, error) {
	rows, _ := r.db.Query(
		ctx,
		searchDoctors,
		q.Specialty,
		q.State,
		q.FacilityID,
		q.Name,
	)

	return pgx.CollectRows(rows, pgx.RowToStructByPos[Doctor])
}
