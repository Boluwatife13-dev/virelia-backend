package records

import "context"

func (r *Repository) ForPatientUser(
	ctx context.Context,
	userID string,
) ([]Entry, error) {
	return r.list(ctx, `e.patient_id = (
		SELECT id FROM patients WHERE user_id = $1
	)`+latestOnly, userID)
}

func (r *Repository) AtFacility(
	ctx context.Context,
	patientID, facilityID string,
) ([]Entry, error) {
	return r.list(ctx,
		"e.patient_id = $1 AND e.facility_id = $2"+latestOnly,
		patientID, facilityID)
}

func (r *Repository) ForAppointment(
	ctx context.Context,
	appointmentID, doctorID string,
) ([]Entry, error) {
	return r.list(ctx, `e.id IN (
		SELECT cr.entry_id
		FROM consultation_records cr
		JOIN appointments a ON a.id = cr.appointment_id
		WHERE a.id = $1 AND a.doctor_id = $2
	)`, appointmentID, doctorID)
}

func (r *Repository) History(
	ctx context.Context,
	id string,
) ([]Entry, error) {
	return r.list(ctx, `e.id IN (
		WITH RECURSIVE chain AS (
			SELECT id, replaces_id FROM record_entries WHERE id = $1
			UNION ALL
			SELECT r.id, r.replaces_id FROM record_entries r
			JOIN chain c ON r.id = c.replaces_id
		)
		SELECT id FROM chain
	)`, id)
}
