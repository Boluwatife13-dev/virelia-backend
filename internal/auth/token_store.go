package auth

import (
	"context"
	"errors"
	"time"
)

func (r *Repository) SaveRefreshToken(
	ctx context.Context,
	t RefreshToken,
	hash, agent, ip string,
	expiresAt time.Time,
) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO refresh_tokens (user_id, family_id,
		  token_hash, user_agent, ip_address, expires_at)
		VALUES ($1, COALESCE(NULLIF($2, '')::uuid, gen_random_uuid()),
		  $3, $4, $5, $6)`,
		t.UserID, t.FamilyID, hash, agent, ip, expiresAt,
	)

	return err
}

func (r *Repository) UseRefreshToken(
	ctx context.Context, hash string,
) (*RefreshToken, error) {
	var t RefreshToken

	err := r.db.QueryRow(ctx, `
		UPDATE refresh_tokens SET revoked_at = now()
		WHERE token_hash = $1 AND revoked_at IS NULL
		RETURNING user_id, family_id, expires_at`, hash,
	).Scan(&t.UserID, &t.FamilyID, &t.ExpiresAt)

	if errors.Is(dbError(err), ErrNotFound) {
		r.db.Exec(ctx, `
			UPDATE refresh_tokens SET revoked_at = now()
			WHERE revoked_at IS NULL AND family_id = (
			  SELECT family_id FROM refresh_tokens
			  WHERE token_hash = $1)`, hash)

		return nil, ErrInvalidToken
	}

	return &t, err
}

func (r *Repository) RevokeTokens(
	ctx context.Context, field, value string,
) error {
	_, err := r.db.Exec(ctx, `
		UPDATE refresh_tokens SET revoked_at = now()
		WHERE `+field+` = $1 AND revoked_at IS NULL`, value,
	)

	return err
}
