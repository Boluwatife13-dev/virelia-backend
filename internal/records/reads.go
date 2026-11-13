package records

import "context"

func (s *Service) MyRecord(ctx context.Context, userID string) ([]Entry, error) {
	return s.repo.ForPatientUser(ctx, userID)
}

func (s *Service) AtMyHospital(
	ctx context.Context,
	d Doctor,
	networkID string,
) ([]Entry, string, error) {
	patientID, err := s.PatientAtMyHospital(ctx, d, networkID)
	if err != nil {
		return nil, "", err
	}

	list, err := s.repo.AtFacility(ctx, patientID, d.FacilityID)

	return list, patientID, err
}

func (s *Service) Attached(
	ctx context.Context,
	d Doctor,
	appointmentID string,
) ([]Entry, error) {
	return s.repo.ForAppointment(ctx, appointmentID, d.ID)
}

func (s *Service) History(
	ctx context.Context,
	d Doctor,
	id string,
) ([]Entry, string, error) {
	e, err := s.repo.get(ctx, id)
	if err != nil {
		return nil, "", err
	}

	if e.FacilityID != d.FacilityID {
		return nil, "", ErrNotAllowed
	}

	if err := s.mustTreat(ctx, d, e.PatientID); err != nil {
		return nil, "", err
	}

	list, err := s.repo.History(ctx, id)

	return list, e.PatientID, err
}
