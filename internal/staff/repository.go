package staff

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrBadInvite = errors.New("this invite is invalid or has expired")
	ErrExists    = errors.New("this email or phone is already registered")
	ErrNotFound  = errors.New("staff member not found")
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// List returns everyone at one hospital.
func (r *Repository) List(
	ctx context.Context,
	facilityID string,
) ([]Member, error) {
	rows, _ := r.db.Query(ctx, `
        SELECT s.user_id, u.email, s.first_name, s.last_name,
               s.role, s.status, s.licence_number, s.specialty
        FROM staff s
        JOIN users u ON u.id = s.user_id
        WHERE s.facility_id = $1
        ORDER BY s.created_at`, facilityID)

	return pgx.CollectRows(rows, pgx.RowToStructByPos[Member])
}

// SetStatus changes a staff member's status,
// but only within the given hospital.
func (r *Repository) SetStatus(
	ctx context.Context,
	facilityID, userID, status string,
) error {
	tag, err := r.db.Exec(ctx, `
        UPDATE staff SET status = $3
        WHERE facility_id = $1 AND user_id = $2`,
		facilityID, userID, status,
	)

	if err == nil && tag.RowsAffected() == 0 {
		return ErrNotFound
	}

	return err
}

// AdminEmails returns a hospital's admin emails.
func (r *Repository) AdminEmails(
	ctx context.Context,
	facilityID string,
) ([]string, error) {
	rows, _ := r.db.Query(ctx, `
        SELECT u.email FROM staff s
        JOIN users u ON u.id = s.user_id
        WHERE s.facility_id = $1
          AND s.role = 'hospital_admin'
          AND s.status = 'active'`, facilityID)

	return pgx.CollectRows(rows, pgx.RowTo[string])
}
