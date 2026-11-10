package patients

import (
	"context"

	"github.com/jackc/pgx/v5"
)

const searchPatients = `
	SELECT
		p.id,
		p.network_id,
		p.first_name,
		p.last_name,
		p.date_of_birth,
		p.sex,
		COALESCE(p.phone, u.phone),
		p.user_id IS NOT NULL
	FROM patients p
	LEFT JOIN users u ON u.id = p.user_id
	WHERE ($1 <> '' AND (p.phone = $1 OR u.phone = $1))
	OR ($2 <> '' AND p.nin_hash = $2)
	OR ($3 <> '' AND p.network_id = $3)
	ORDER BY p.created_at DESC
	LIMIT 20`

func (r *Repository) Search(
	ctx context.Context,
	phone string,
	ninHash string,
	networkID string,
) ([]Patient, error) {
	rows, _ := r.db.Query(
		ctx,
		searchPatients,
		phone,
		ninHash,
		networkID,
	)

	return pgx.CollectRows(rows, pgx.RowToStructByPos[Patient])
}
