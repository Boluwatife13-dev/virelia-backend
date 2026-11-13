package directory

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func (r *Repository) SaveProfile(
	ctx context.Context,
	doctorID string,
	p ProfileRequest,
) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO doctor_profiles
			(user_id, fee_naira, slot_minutes, bio, accepting)
		VALUES ($1, $2, $3, NULLIF($4, ''), $5)
		ON CONFLICT (user_id) DO UPDATE SET
			fee_naira = EXCLUDED.fee_naira,
			slot_minutes = EXCLUDED.slot_minutes,
			bio = EXCLUDED.bio,
			accepting = EXCLUDED.accepting,
			updated_at = now()`,
		doctorID, p.FeeNaira, p.SlotMinutes, p.Bio, p.Accepting)

	return err
}

func (r *Repository) SetHours(
	ctx context.Context,
	doctorID string,
	hours []Hours,
) error {
	return pgx.BeginFunc(ctx, r.db, func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx,
			`DELETE FROM doctor_hours WHERE doctor_id = $1`,
			doctorID)
		if err != nil {
			return err
		}

		for _, h := range hours {
			_, err = tx.Exec(ctx, `
				INSERT INTO doctor_hours
					(doctor_id, weekday, starts, ends)
				VALUES ($1, $2, $3::time, $4::time)`,
				doctorID, h.Weekday, h.Starts, h.Ends)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *Repository) GetHours(
	ctx context.Context,
	doctorID string,
) ([]Hours, error) {
	rows, _ := r.db.Query(ctx, `
		SELECT weekday, to_char(starts, 'HH24:MI'), to_char(ends, 'HH24:MI')
		FROM doctor_hours
		WHERE doctor_id = $1
		ORDER BY weekday, starts`,
		doctorID)

	return pgx.CollectRows(rows, pgx.RowToStructByPos[Hours])
}
