package auth

import (
	"context"
	"errors"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"virelia/internal/platform/crypt"
)

var (
	ErrWrongLogin  = errors.New("invalid email or password")
	ErrNotVerified = errors.New("please verify your email first")
	ErrInactive    = errors.New("this account is not active")
	ErrBadDOB      = errors.New("enter a valid date of birth")
)

type Mailer interface {
	Send(ctx context.Context, to, subject, html string) error
}

type Service struct {
	repo   *Repository
	mailer Mailer
	secret string
	box    *crypt.Box
}

func NewService(
	repo *Repository,
	m Mailer,
	secret string,
	box *crypt.Box,
) *Service {
	return &Service{
		repo:   repo,
		mailer: m,
		secret: secret,
		box:    box,
	}
}

func (s *Service) Register(ctx context.Context, req RegisterRequest) error {
	dob, err := time.Parse("2006-01-02", req.DateOfBirth)
	if err != nil || dob.After(time.Now()) {
		return ErrBadDOB
	}

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password), bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	user := &User{
		Email:        normalize(req.Email),
		Phone:        req.Phone,
		PasswordHash: string(hash),
	}

	patient := &Patient{
		NetworkID:      newNetworkID(),
		FirstName:      strings.TrimSpace(req.FirstName),
		LastName:       strings.TrimSpace(req.LastName),
		DateOfBirth:    dob,
		Sex:            req.Sex,
		State:          optional(req.State),
		Address:        optional(req.Address),
		EmergencyName:  optional(req.EmergencyName),
		EmergencyPhone: optional(req.EmergencyPhone),
	}

	if req.NIN != "" {
		ninHash := hashCode(req.NIN, s.secret)
		patient.NINHash = &ninHash
	}

	if err := s.repo.CreatePatient(ctx, user, patient); err != nil {
		return err
	}

	return s.sendCode(ctx, user, PurposeVerifyEmail)
}
