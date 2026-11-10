package guardians

import (
	"context"

	"github.com/jackc/pgx/v5"
)

const linkColumns = `
    SELECT
        l.id,
        p.network_id,
        p.first_name,
        p.last_name,
        l.relationship,
        l.verified_at IS NOT NULL,
        l.expires_at
    FROM guardians l`

func (r *Repository) MyChildren(
	ctx context.Context,
	userID string,
) ([]Link, error) {
	rows, _ := r.db.Query(ctx, linkColumns+`
        JOIN patients p ON p.id = l.child_id
        JOIN patients g ON g.id = l.guardian_id
        WHERE g.user_id = $1
        AND l.revoked_at IS NULL
        AND l.expires_at > now()
        ORDER BY p.date_of_birth`, userID)

	return pgx.CollectRows(rows, pgx.RowToStructByPos[Link])
}

func (r *Repository) GuardiansOf(
	ctx context.Context,
	childNetworkID string,
) ([]Link, error) {
	rows, _ := r.db.Query(ctx, linkColumns+`
        JOIN patients p ON p.id = l.guardian_id
        JOIN patients c ON c.id = l.child_id
        WHERE c.network_id = $1
        AND l.revoked_at IS NULL
        ORDER BY l.created_at`, childNetworkID)

	return pgx.CollectRows(rows, pgx.RowToStructByPos[Link])
}

func (r *Repository) update(
	ctx context.Context,
	query string,
	args ...any,
) error {
	tag, err := r.db.Exec(ctx, query, args...)

	if err == nil && tag.RowsAffected() == 0 {
		return ErrNotFound
	}

	return err
}

func (r *Repository) Verify(ctx context.Context, linkID string) error {
	return r.update(ctx, `
        UPDATE guardians SET verified_at = now()
        WHERE id = $1
        AND verified_at IS NULL
        AND revoked_at IS NULL`, linkID)
}

func (r *Repository) Revoke(ctx context.Context, linkID string) error {
	return r.update(ctx, `
        UPDATE guardians SET revoked_at = now()
        WHERE id = $1
        AND revoked_at IS NULL`, linkID)
}

func (r *Repository) Leave(
	ctx context.Context,
	linkID string,
	userID string,
) error {
	return r.update(ctx, `
        UPDATE guardians SET revoked_at = now()
        WHERE id = $1
        AND revoked_at IS NULL
        AND guardian_id = (
            SELECT id FROM patients WHERE user_id = $2
        )`, linkID, userID)
}
