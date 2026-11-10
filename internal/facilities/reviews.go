package facilities

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

var ErrNeedsSecond = errors.New(
	"a different admin must check this hospital first",
)

func (r *Repository) List(
	ctx context.Context,
	status string,
) ([]Facility, error) {
	rows, _ := r.db.Query(ctx, `
        SELECT `+columns+`
        FROM facilities
        WHERE $1 = '' OR status = $1
        ORDER BY created_at DESC`, status)

	return pgx.CollectRows(rows, toFacility)
}

func (r *Repository) update(
	ctx context.Context,
	notFound error,
	query string,
	args ...any,
) (*Facility, error) {
	var f Facility

	row := r.db.QueryRow(ctx, query+" RETURNING "+columns, args...)

	err := scan(row, &f)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, notFound
	}

	return &f, err
}
