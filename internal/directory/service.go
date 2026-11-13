package directory

import (
	"context"
	"errors"
	"strings"
)

var ErrBadHours = errors.New("each block must end after it starts")

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo}
}

func (s *Service) Search(ctx context.Context, q Search) ([]Doctor, error) {
	q.Specialty = strings.TrimSpace(q.Specialty)
	q.Name = strings.TrimSpace(q.Name)

	return s.repo.Search(ctx, q)
}

func (s *Service) SaveProfile(
	ctx context.Context,
	doctorID string,
	p ProfileRequest,
) error {
	return s.repo.SaveProfile(ctx, doctorID, p)
}

func (s *Service) SetHours(
	ctx context.Context,
	doctorID string,
	hours []Hours,
) error {
	for _, h := range hours {
		if h.Starts >= h.Ends {
			return ErrBadHours
		}
	}

	return s.repo.SetHours(ctx, doctorID, hours)
}

func (s *Service) Hours(ctx context.Context, doctorID string) ([]Hours, error) {
	return s.repo.GetHours(ctx, doctorID)
}
