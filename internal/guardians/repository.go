package guardians

import (
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrNoGuardian = errors.New("guardian not found")
	ErrExists     = errors.New("a patient with this NIN already exists")
	ErrNotFound   = errors.New("link not found")
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db}
}

type NewChild struct {
	NetworkID    string
	FirstName    string
	LastName     string
	DateOfBirth  time.Time
	Sex          string
	NINHash      string
	Relationship string
	AddedBy      string
	FacilityID   string
	Verified     bool
}
