package appointments

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrNoDoctor  = errors.New("doctor not found or not taking bookings")
	ErrNoPatient = errors.New("patient not found")
	ErrTaken     = errors.New("that time is no longer free")
	ErrNotFound  = errors.New("appointment not found or cannot be changed")
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db}
}

type doctor struct {
	FacilityID  string
	SlotMinutes int
}

func (r *Repository) doctor(ctx context.Context, id string) (doctor, error) {
	var d doctor

	err := r.db.QueryRow(ctx, `
		SELECT s.facility_id, p.slot_minutes
		FROM staff s
		JOIN doctor_profiles p ON p.user_id = s.user_id
		JOIN facilities f ON f.id = s.facility_id
		WHERE s.user_id = $1
		AND s.role = 'doctor'
		AND s.status = 'active'
		AND f.status = 'approved'
		AND p.accepting`, id,
	).Scan(&d.FacilityID, &d.SlotMinutes)

	if errors.Is(err, pgx.ErrNoRows) {
		return d, ErrNoDoctor
	}

	return d, err
}

func (r *Repository) patientFor(
	ctx context.Context,
	column, value string,
) (string, error) {
	var id string

	err := r.db.QueryRow(ctx,
		`SELECT id FROM patients WHERE `+column+` = $1`,
		value,
	).Scan(&id)

	if errors.Is(err, pgx.ErrNoRows) {
		return "", ErrNoPatient
	}

	return id, err
}

type newAppointment struct {
	PatientID  string
	DoctorID   string
	FacilityID string
	Slot       Slot
	Reason     string
	BookedBy   string
}

func (r *Repository) insert(
	ctx context.Context,
	a newAppointment,
) (string, error) {
	var id string

	err := r.db.QueryRow(ctx, `
		INSERT INTO appointments (
			patient_id, doctor_id, facility_id,
			starts_at, ends_at, reason, booked_by
		)
		VALUES ($1, $2, $3, $4, $5, NULLIF($6, ''), $7)
		RETURNING id`,
		a.PatientID, a.DoctorID, a.FacilityID,
		a.Slot.StartsAt, a.Slot.EndsAt, a.Reason, a.BookedBy,
	).Scan(&id)

	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) && pgErr.Code == "23P01" {
		return "", ErrTaken
	}

	return id, err
}
