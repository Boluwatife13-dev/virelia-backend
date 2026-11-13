package appointments

import "context"

func (s *Service) MyAppointments(
	ctx context.Context,
	userID string,
) ([]Appointment, error) {
	return s.repo.forPatientUser(ctx, userID)
}

func (s *Service) Day(
	ctx context.Context,
	column, value, date string,
) ([]Appointment, error) {
	day, err := parseDay(date)
	if err != nil {
		return nil, err
	}

	return s.repo.onDay(ctx, column, value, day)
}

func (s *Service) CancelMine(
	ctx context.Context,
	userID, id, reason string,
) error {
	patientID, err := s.repo.patientFor(ctx, "user_id", userID)
	if err != nil {
		return err
	}

	return s.repo.setStatus(ctx, change{
		ID:          id,
		To:          "cancelled",
		From:        []string{"booked"},
		ScopeColumn: "patient_id",
		ScopeValue:  patientID,
		ByUser:      userID,
		Reason:      reason,
	})
}

func (s *Service) AtDesk(
	ctx context.Context,
	facilityID, staffID, id, to, reason string,
) error {
	from := []string{"booked"}

	if to == "cancelled" {
		from = []string{"booked", "checked_in"}
	}

	return s.repo.setStatus(ctx, change{
		ID:          id,
		To:          to,
		From:        from,
		ScopeColumn: "facility_id",
		ScopeValue:  facilityID,
		ByUser:      staffID,
		Reason:      reason,
	})
}

func (s *Service) Complete(ctx context.Context, doctorID, id string) error {
	return s.repo.setStatus(ctx, change{
		ID:          id,
		To:          "completed",
		From:        []string{"booked", "checked_in"},
		ScopeColumn: "doctor_id",
		ScopeValue:  doctorID,
		ByUser:      doctorID,
	})
}
