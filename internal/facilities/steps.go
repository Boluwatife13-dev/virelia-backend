package facilities

import "context"

func (r *Repository) Check(
	ctx context.Context,
	id, adminID, note string,
) (*Facility, error) {
	return r.update(ctx, ErrNotPending, `
        UPDATE facilities SET
            status = 'checked',
            checked_by = $2,
            checked_at = now(),
            check_note = $3,
            updated_at = now()
        WHERE id = $1
        AND status = 'pending'`,
		id, adminID, note)
}

func (r *Repository) Approve(
	ctx context.Context,
	id, adminID string,
) (*Facility, error) {
	return r.update(ctx, ErrNeedsSecond, `
        UPDATE facilities SET
            status = 'approved',
            reviewed_by = $2,
            reviewed_at = now(),
            updated_at = now()
        WHERE id = $1
        AND status = 'checked'
        AND checked_by <> $2`,
		id, adminID)
}

func (r *Repository) Reject(
	ctx context.Context,
	id, adminID, reason string,
) (*Facility, error) {
	return r.update(ctx, ErrNotPending, `
        UPDATE facilities SET
            status = 'rejected',
            rejection_reason = $3,
            reviewed_by = $2,
            reviewed_at = now(),
            updated_at = now()
        WHERE id = $1
        AND status IN ('pending', 'checked')`,
		id, adminID, reason)
}
