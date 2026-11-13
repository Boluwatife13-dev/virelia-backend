package appointments

import "context"

type change struct {
	ID          string
	To          string
	From        []string
	ScopeColumn string
	ScopeValue  string
	ByUser      string
	Reason      string
}

func (r *Repository) setStatus(ctx context.Context, c change) error {
	tag, err := r.db.Exec(ctx, `
		UPDATE appointments SET
			status = $2,
			cancelled_by = CASE WHEN $2 = 'cancelled'
				THEN $5::uuid END,
			cancel_reason = CASE WHEN $2 = 'cancelled'
				THEN NULLIF($6, '') END,
			updated_at = now()
		WHERE id = $1
		AND status = ANY($3)
		AND `+c.ScopeColumn+` = $4`,
		c.ID, c.To, c.From, c.ScopeValue, c.ByUser, c.Reason)

	if err == nil && tag.RowsAffected() == 0 {
		return ErrNotFound
	}

	return err
}
