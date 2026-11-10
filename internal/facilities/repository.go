package facilities

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotPending = errors.New(
	"hospital not found or already reviewed",
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db}
}

const insertFacility = `
	INSERT INTO facilities (
		name,
		facility_type,
		licence_number,
		licence_expires_at,
		state,
		address,
		email,
		phone,
		admin_first_name,
		admin_last_name,
		admin_email
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	RETURNING id, status, created_at`

func (r *Repository) Create(
	ctx context.Context,
	f *Facility,
) error {
	row := r.db.QueryRow(
		ctx,
		insertFacility,
		f.Name,
		f.Type,
		f.LicenceNumber,
		f.LicenceExpiresAt,
		f.State,
		f.Address,
		f.Email,
		f.Phone,
		f.AdminFirstName,
		f.AdminLastName,
		f.AdminEmail,
	)

	return row.Scan(
		&f.ID,
		&f.Status,
		&f.CreatedAt,
	)
}

const adminEmails = `
	SELECT u.email
	FROM staff s
	JOIN users u ON u.id = s.user_id
	WHERE s.role = 'platform_admin'
	AND s.status = 'active'`

func (r *Repository) PlatformAdminEmails(
	ctx context.Context,
) ([]string, error) {
	rows, _ := r.db.Query(ctx, adminEmails)

	return pgx.CollectRows(rows, pgx.RowTo[string])
}
