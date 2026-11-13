package records

import "context"

func (s *Service) Add(
	ctx context.Context,
	d Doctor,
	req AddRequest,
) (string, string, error) {
	if err := checkSize(req.Data); err != nil {
		return "", "", err
	}

	patientID, err := s.PatientAtMyHospital(ctx, d, req.PatientNetworkID)
	if err != nil {
		return "", "", err
	}

	id, err := s.repo.insert(ctx, newEntry{
		PatientID:     patientID,
		FacilityID:    d.FacilityID,
		AuthorID:      d.ID,
		AppointmentID: req.AppointmentID,
		Kind:          req.Kind,
		Data:          req.Data,
		Version:       1,
	})

	return id, patientID, err
}

func (s *Service) Correct(
	ctx context.Context,
	d Doctor,
	id string,
	req CorrectRequest,
) (string, string, error) {
	if err := checkSize(req.Data); err != nil {
		return "", "", err
	}

	old, err := s.repo.get(ctx, id)
	if err != nil {
		return "", "", err
	}

	if old.FacilityID != d.FacilityID {
		return "", "", ErrNotAllowed
	}

	if err := s.mustTreat(ctx, d, old.PatientID); err != nil {
		return "", "", err
	}

	if latest, err := s.repo.isLatest(ctx, id); err != nil || !latest {
		return "", "", ErrOldVersion
	}

	newID, err := s.repo.insert(ctx, newEntry{
		PatientID:  old.PatientID,
		FacilityID: old.FacilityID,
		AuthorID:   d.ID,
		Kind:       old.Kind,
		Data:       req.Data,
		Version:    old.Version + 1,
		ReplacesID: old.ID,
		Reason:     req.Reason,
	})

	return newID, old.PatientID, err
}
