package facilities

import "github.com/jackc/pgx/v5"

const columns = `
    id, name, facility_type, licence_number,
    licence_expires_at, state, address, email, phone,
    admin_first_name, admin_last_name, admin_email,
    status, rejection_reason, checked_by,
    check_note, created_at`

func scan(row pgx.Row, f *Facility) error {
	return row.Scan(
		&f.ID, &f.Name, &f.Type, &f.LicenceNumber,
		&f.LicenceExpiresAt, &f.State, &f.Address,
		&f.Email, &f.Phone, &f.AdminFirstName,
		&f.AdminLastName, &f.AdminEmail, &f.Status,
		&f.RejectionReason, &f.CheckedBy,
		&f.CheckNote, &f.CreatedAt,
	)
}

func toFacility(row pgx.CollectableRow) (Facility, error) {
	var f Facility
	err := scan(row, &f)

	return f, err
}
