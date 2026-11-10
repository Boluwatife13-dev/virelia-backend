package patients

import (
	"context"
	"errors"
	"strings"
	"time"
)

var (
	ErrBadBirthDate = errors.New("date of birth can't be in the future")
	ErrNoSearch     = errors.New("give a phone, NIN or network ID")
)

type IDs interface {
	HashNIN(nin string) string
	NewNetworkID() string
}

type Service struct {
	repo *Repository
	ids  IDs
}

func NewService(repo *Repository, ids IDs) *Service {
	return &Service{repo, ids}
}

func (s *Service) Register(
	ctx context.Context,
	req RegisterRequest,
	staffID string,
	facilityID string,
) (string, error) {
	dob, err := time.Parse("2006-01-02", req.DateOfBirth)

	if err != nil || dob.After(time.Now()) {
		return "", ErrBadBirthDate
	}

	p := NewPatient{
		NetworkID:      s.ids.NewNetworkID(),
		FirstName:      strings.TrimSpace(req.FirstName),
		LastName:       strings.TrimSpace(req.LastName),
		DateOfBirth:    dob,
		Sex:            req.Sex,
		Phone:          req.Phone,
		State:          req.State,
		Address:        req.Address,
		EmergencyName:  req.EmergencyName,
		EmergencyPhone: req.EmergencyPhone,
		RegisteredBy:   staffID,
		FacilityID:     facilityID,
	}

	if req.NIN != "" {
		p.NINHash = s.ids.HashNIN(req.NIN)
	}

	err = s.repo.Create(ctx, p)

	return p.NetworkID, err
}

func (s *Service) Find(
	ctx context.Context,
	q Search,
) ([]Patient, error) {
	if q.Phone == "" && q.NIN == "" && q.NetworkID == "" {
		return nil, ErrNoSearch
	}

	ninHash := ""

	if q.NIN != "" {
		ninHash = s.ids.HashNIN(q.NIN)
	}

	networkID := strings.ToUpper(strings.TrimSpace(q.NetworkID))

	return s.repo.Search(ctx, q.Phone, ninHash, networkID)
}
