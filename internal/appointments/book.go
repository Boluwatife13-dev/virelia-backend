package appointments

import (
	"context"
	"strings"
	"time"
)

func (s *Service) book(
	ctx context.Context,
	patientID, bookedBy string,
	req BookRequest,
	onlyFacility string,
) (string, error) {
	d, err := s.repo.doctor(ctx, req.DoctorID)
	if err != nil {
		return "", err
	}

	if onlyFacility != "" && d.FacilityID != onlyFacility {
		return "", ErrOtherHospital
	}

	length := time.Duration(d.SlotMinutes) * time.Minute
	day := req.StartsAt.In(wat)

	free, err := s.repo.freeSlots(ctx, req.DoctorID, day, length)
	if err != nil {
		return "", err
	}

	for _, slot := range free {
		if slot.StartsAt.Equal(req.StartsAt) {
			return s.repo.insert(ctx, newAppointment{
				PatientID:  patientID,
				DoctorID:   req.DoctorID,
				FacilityID: d.FacilityID,
				Slot:       slot,
				Reason:     strings.TrimSpace(req.Reason),
				BookedBy:   bookedBy,
			})
		}
	}

	return "", ErrBadTime
}

func (s *Service) BookForMe(
	ctx context.Context,
	userID string,
	req BookRequest,
) (string, error) {
	patientID, err := s.repo.patientFor(ctx, "user_id", userID)
	if err != nil {
		return "", err
	}

	return s.book(ctx, patientID, userID, req, "")
}

func (s *Service) BookAtDesk(
	ctx context.Context,
	staffID, facilityID string,
	req DeskBookRequest,
) (string, error) {
	networkID := strings.ToUpper(strings.TrimSpace(req.PatientNetworkID))

	patientID, err := s.repo.patientFor(ctx, "network_id", networkID)
	if err != nil {
		return "", err
	}

	return s.book(ctx, patientID, staffID, req.BookRequest, facilityID)
}
