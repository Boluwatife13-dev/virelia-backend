package keys

import (
	"context"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db}
}

// Admin loads a user's email and saved keys.
func (r *Repository) Admin(
	ctx context.Context,
	userID string,
) (*admin, error) {
	a := &admin{id: userID}

	err := r.db.QueryRow(ctx,
		`SELECT email FROM users WHERE id = $1`,
		userID,
	).Scan(&a.email)
	if err != nil {
		return nil, err
	}

	rows, _ := r.db.Query(ctx,
		`SELECT credential FROM security_keys WHERE user_id = $1`,
		userID,
	)

	a.keys, err = pgx.CollectRows(
		rows,
		pgx.RowTo[webauthn.Credential],
	)

	return a, err
}

// SaveKey stores a newly registered key.
func (r *Repository) SaveKey(
	ctx context.Context,
	userID string,
	cred *webauthn.Credential,
) error {
	_, err := r.db.Exec(ctx, `
        INSERT INTO security_keys
            (user_id, credential_id, credential)
        VALUES ($1, $2, $3)`,
		userID, cred.ID, cred)

	return err
}

// UpdateKey saves the new use counter after a login.
func (r *Repository) UpdateKey(
	ctx context.Context,
	cred *webauthn.Credential,
) error {
	_, err := r.db.Exec(ctx, `
        UPDATE security_keys
        SET credential = $2,
            last_used_at = now()
        WHERE credential_id = $1`,
		cred.ID, cred)

	return err
}
