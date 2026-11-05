package auth

import (
	"context"
	"errors"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"

	"virelia-backend/models"
	repo "virelia-backend/repositories/auth"
	"virelia-backend/utils"
)

var (
	ErrInvalidLogin    = errors.New("invalid email or password")
	ErrAccountInactive = errors.New("account is not active")
	ErrInvalidToken    = errors.New("invalid refresh token")
)

type Service struct {
	users     *repo.UserRepository
	tokens    *repo.RefreshTokenRepository
	otps      *repo.OTPRepository
	mailer    utils.Mailer
	jwtSecret string
}

func NewService(
	users *repo.UserRepository,
	tokens *repo.RefreshTokenRepository,
	otps *repo.OTPRepository,
	mailer utils.Mailer,
	jwtSecret string,
) *Service {
	return &Service{users, tokens, otps, mailer, jwtSecret}
}

type RegisterInput struct {
	FirstName, LastName, Email, Phone, Password string
}

func (s *Service) RegisterPatient(
	ctx context.Context, in RegisterInput,
) (*models.User, error) {
	hash, err := utils.HashPassword(in.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		FirstName:    strings.TrimSpace(in.FirstName),
		LastName:     strings.TrimSpace(in.LastName),
		Email:        normalizeEmail(in.Email),
		Phone:        strings.TrimSpace(in.Phone),
		PasswordHash: hash,
		Role:         models.RolePatient,
		Status:       models.StatusActive,
	}

	return user, s.users.Create(ctx, user)
}

func (s *Service) Login(
	ctx context.Context,
	email, password, userAgent, ip string,
) (*TokenPair, error) {
	user, err := s.users.FindByEmail(ctx, normalizeEmail(email))
	if errors.Is(err, repo.ErrUserNotFound) {
		return nil, ErrInvalidLogin
	}
	if err != nil {
		return nil, err
	}

	if !utils.CheckPassword(password, user.PasswordHash) {
		return nil, ErrInvalidLogin
	}

	if user.Status != models.StatusActive {
		return nil, ErrAccountInactive
	}

	err = s.users.RecordSuccessfulLogin(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return s.issueTokens(ctx, user, "", userAgent, ip)
}

func (s *Service) GetUser(
	ctx context.Context, id string,
) (*models.User, error) {
	userID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	return s.users.FindByID(ctx, userID)
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}