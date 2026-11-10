package patients

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrExists = errors.New("a patient with this NIN already exists")

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db}
}

// NewPatient is everything saved for a walk-in.
type NewPatient struct {
	NetworkID      string
	FirstName      string
	LastName       string
	DateOfBirth    time.Time
	Sex            string
	Phone          string
	NINHash        string
	State          string
	Address        string
	EmergencyName  string
	EmergencyPhone string
	RegisteredBy   string
	FacilityID     string
}

const insertPatient = `
    INSERT INTO patients (
        network_id,
        first_name,
        last_name,
        date_of_birth,
        sex,
        phone,
        nin_hash,
        state,
        address,
        emergency_contact_name,
        emergency_contact_phone,
        registered_by,
        registered_facility_id
    )
    VALUES (
        $1, $2, $3, $4, $5,
        NULLIF($6, ''),
        NULLIF($7, ''),
        NULLIF($8, ''),
        NULLIF($9, ''),
        NULLIF($10, ''),
        NULLIF($11, ''),
        $12, $13
    )`

// Create saves a walk-in patient.
func (r *Repository) Create(
	ctx context.Context,
	p NewPatient,
) error {
	_, err := r.db.Exec(
		ctx,
		insertPatient,
		p.NetworkID,
		p.FirstName,
		p.LastName,
		p.DateOfBirth,
		p.Sex,
		p.Phone,
		p.NINHash,
		p.State,
		p.Address,
		p.EmergencyName,
		p.EmergencyPhone,
		p.RegisteredBy,
		p.FacilityID,
	)

	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return ErrExists
	}

	return err
}
