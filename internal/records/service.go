package records

import (
	"context"
	"encoding/json"
	"strings"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo}
}

type Doctor struct {
	ID         string
	FacilityID string
}

func checkSize(data map[string]any) error {
	b, err := json.Marshal(data)
	if err != nil || len(b) > 20000 || len(data) == 0 {
		return ErrTooBig
	}

	return nil
}

func (s *Service) mustTreat(
	ctx context.Context,
	d Doctor,
	patientID string,
) error {
	ok, err := s.repo.treats(ctx, d.ID, d.FacilityID, patientID)
	if err != nil {
		return err
	}

	if !ok {
		return ErrNotAllowed
	}

	return nil
}

func (s *Service) PatientAtMyHospital(
	ctx context.Context,
	d Doctor,
	networkID string,
) (string, error) {
	networkID = strings.ToUpper(strings.TrimSpace(networkID))

	patientID, err := s.repo.patientByNetworkID(ctx, networkID)
	if err != nil {
		return "", err
	}

	return patientID, s.mustTreat(ctx, d, patientID)
}
