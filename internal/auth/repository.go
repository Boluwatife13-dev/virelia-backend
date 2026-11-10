package auth

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrNotFound = errors.New("not found")
	ErrExists   = errors.New("email, phone or NIN already registered")
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreatePatient(
	ctx context.Context, u *User, p *Patient,
) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	err = tx.QueryRow(ctx, `
		INSERT INTO users (email, phone, password_hash)
		VALUES ($1, $2, $3) RETURNING id`,
		u.Email, u.Phone, u.PasswordHash,
	).Scan(&u.ID)
	if err != nil {
		return dbError(err)
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO patients (
		  user_id, network_id, first_name, last_name,
		  date_of_birth, sex, nin_hash, state, address,
		  emergency_contact_name, emergency_contact_phone,
		  terms_accepted_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,now())`,
		u.ID, p.NetworkID, p.FirstName, p.LastName,
		p.DateOfBirth, p.Sex, p.NINHash, p.State,
		p.Address, p.EmergencyName, p.EmergencyPhone,
	)
	if err != nil {
		return dbError(err)
	}

	return tx.Commit(ctx)
}

func (r *Repository) FindUser(
	ctx context.Context, field, value string,
) (*User, error) {
	var u User

	err := r.db.QueryRow(ctx, `
		SELECT id, email, phone, password_hash,
		       status, email_verified_at
		FROM users WHERE `+field+` = $1`, value,
	).Scan(
		&u.ID, &u.Email, &u.Phone, &u.PasswordHash,
		&u.Status, &u.EmailVerifiedAt,
	)

	return &u, dbError(err)
}

func (r *Repository) FindPatient(
	ctx context.Context, userID string,
) (*Patient, error) {
	var p Patient

	err := r.db.QueryRow(ctx, `
		SELECT network_id, first_name, last_name,
		       date_of_birth, sex, state, address
		FROM patients WHERE user_id = $1`, userID,
	).Scan(
		&p.NetworkID, &p.FirstName, &p.LastName,
		&p.DateOfBirth, &p.Sex, &p.State, &p.Address,
	)

	return &p, dbError(err)
}

func (r *Repository) UpdateUser(
	ctx context.Context, id, set string, args ...any,
) error {
	_, err := r.db.Exec(ctx,
		"UPDATE users SET "+set+", updated_at = now() WHERE id = $1",
		append([]any{id}, args...)...,
	)

	return err
}

func dbError(err error) error {
	var pgErr *pgconn.PgError

	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return ErrNotFound
	case errors.As(err, &pgErr) && pgErr.Code == "23505":
		return ErrExists
	}

	return err
}
