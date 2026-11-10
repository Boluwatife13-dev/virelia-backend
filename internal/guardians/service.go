package guardians

import (
	"errors"
	"strings"
	"time"
)

var ErrNotChild = errors.New("the child must be under 18")

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

func (s *Service) newChild(c Child) (NewChild, error) {
	dob, err := time.Parse("2006-01-02", c.DateOfBirth)

	adult := time.Now().AddDate(-18, 0, 0)

	if err != nil || dob.After(time.Now()) || !dob.After(adult) {
		return NewChild{}, ErrNotChild
	}

	n := NewChild{
		NetworkID:    s.ids.NewNetworkID(),
		FirstName:    strings.TrimSpace(c.FirstName),
		LastName:     strings.TrimSpace(c.LastName),
		DateOfBirth:  dob,
		Sex:          c.Sex,
		Relationship: c.Relationship,
	}

	if c.NIN != "" {
		n.NINHash = s.ids.HashNIN(c.NIN)
	}

	return n, nil
}
