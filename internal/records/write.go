package records

import "context"

type newEntry struct {
	PatientID     string
	FacilityID    string
	AuthorID      string
	AppointmentID string
	Kind          string
	Data          map[string]any
	Version       int
	ReplacesID    string
	Reason        string
}

func (r *Repository) insert(ctx context.Context, e newEntry) (string, error) {
	var id string

	err := r.db.QueryRow(ctx, `
		INSERT INTO record_entries (
			patient_id, facility_id, author_id,
			appointment_id, kind, data, version,
			replaces_id, correction_reason
		)
		VALUES (
			$1, $2, $3,
			NULLIF($4, '')::uuid, $5, $6, $7,
			NULLIF($8, '')::uuid, NULLIF($9, '')
		)
		RETURNING id`,
		e.PatientID, e.FacilityID, e.AuthorID,
		e.AppointmentID, e.Kind, e.Data, e.Version,
		e.ReplacesID, e.Reason,
	).Scan(&id)

	return id, err
}
