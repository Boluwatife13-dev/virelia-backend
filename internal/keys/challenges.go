package keys

import (
	"context"
	"errors"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/jackc/pgx/v5"
)

var ErrExpired = errors.New("this request expired, please start again")

// SaveChallenge keeps the challenge for 5 minutes,
// replacing any older one for the same purpose.
func (r *Repository) SaveChallenge(
	ctx context.Context,
	userID string,
	purpose string,
	data *webauthn.SessionData,
) error {
	_, err := r.db.Exec(ctx, `
        INSERT INTO key_challenges
            (user_id, purpose, data, expires_at)
        VALUES ($1, $2, $3, now() + interval '5 minutes')
        ON CONFLICT (user_id, purpose) DO UPDATE
        SET data = EXCLUDED.data,
            expires_at = EXCLUDED.expires_at`,
		userID, purpose, data)

	return err
}

// TakeChallenge reads and deletes the challenge,
// so it can only be used once.
func (r *Repository) TakeChallenge(
	ctx context.Context,
	userID string,
	purpose string,
) (webauthn.SessionData, error) {
	var data webauthn.SessionData

	err := r.db.QueryRow(ctx, `
        DELETE FROM key_challenges
        WHERE user_id = $1
        AND purpose = $2
        AND expires_at > now()
        RETURNING data`,
		userID, purpose,
	).Scan(&data)

	if errors.Is(err, pgx.ErrNoRows) {
		return data, ErrExpired
	}

	return data, err
}
