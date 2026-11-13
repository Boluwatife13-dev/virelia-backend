package appointments

import (
	"context"
	"errors"
	"time"
)

var (
	ErrBadDate       = errors.New("date must look like 2026-12-01")
	ErrBadTime       = errors.New("pick one of the doctor's free times")
	ErrOtherHospital = errors.New("this doctor works at another hospital")
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo}
}

func parseDay(date string) (time.Time, error) {
	if date == "" {
		return time.Now().In(wat), nil
	}

	day, err := time.ParseInLocation("2006-01-02", date, wat)
	if err != nil {
		return day, ErrBadDate
	}

	return day, nil
}

func (s *Service) Slots(
	ctx context.Context,
	doctorID, date string,
) ([]Slot, error) {
	day, err := parseDay(date)
	if err != nil {
		return nil, err
	}

	d, err := s.repo.doctor(ctx, doctorID)
	if err != nil {
		return nil, err
	}

	length := time.Duration(d.SlotMinutes) * time.Minute

	return s.repo.freeSlots(ctx, doctorID, day, length)
}
