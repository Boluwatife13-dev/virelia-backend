package guardians

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
)

const addChild = `
    WITH g AS (
        SELECT id FROM patients
        WHERE %s = $1
    ),
    c AS (
        INSERT INTO patients (
            network_id,
            first_name,
            last_name,
            date_of_birth,
            sex,
            nin_hash,
            registered_by,
            registered_facility_id
        )
        SELECT $2, $3, $4, $5, $6,
            NULLIF($7, ''),
            $8,
            NULLIF($9, '')::uuid
        FROM g
        RETURNING id, date_of_birth
    )
    INSERT INTO guardians (
        child_id,
        guardian_id,
        relationship,
        added_by,
        verified_at,
        expires_at
    )
    SELECT
        c.id,
        g.id,
        $10,
        $8,
        CASE WHEN $11 THEN now() END,
        c.date_of_birth + interval '18 years'
    FROM c, g`

func (r *Repository) AddChild(
	ctx context.Context,
	findBy string,
	guardian string,
	c NewChild,
) error {
	query := fmt.Sprintf(addChild, findBy)

	tag, err := r.db.Exec(
		ctx,
		query,
		guardian,
		c.NetworkID,
		c.FirstName,
		c.LastName,
		c.DateOfBirth,
		c.Sex,
		c.NINHash,
		c.AddedBy,
		c.FacilityID,
		c.Relationship,
		c.Verified,
	)

	var pgErr *pgconn.PgError

	switch {
	case errors.As(err, &pgErr) && pgErr.Code == "23505":
		return ErrExists
	case err == nil && tag.RowsAffected() == 0:
		return ErrNoGuardian
	}

	return err
}
