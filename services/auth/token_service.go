package auth

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"

	"virelia-backend/models"
	repo "virelia-backend/repositories/auth"
	"virelia-backend/utils"
)

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (s *Service) Refresh(
	ctx context.Context, token, userAgent, ip string,
) (*TokenPair, error) {
	record, err := s.findToken(ctx, token)
	if err != nil {
		return nil, err
	}

	err = s.tokens.Revoke(ctx, record.ID)
	if errors.Is(err, repo.ErrTokenRevoked) {
		err = s.tokens.RevokeFamily(ctx, record.FamilyID)
		return nil, errors.Join(ErrInvalidToken, err)
	}
	if err != nil {
		return nil, err
	}

	user, err := s.users.FindByID(ctx, record.UserID)
	if err != nil {
		return nil, err
	}
	if user.Status != models.StatusActive {
		return nil, ErrAccountInactive
	}

	return s.issueTokens(ctx, user, record.FamilyID, userAgent, ip)
}

func (s *Service) Logout(ctx context.Context, token string) error {
	record, err := s.findToken(ctx, token)
	if errors.Is(err, ErrInvalidToken) {
		return nil
	}
	if err != nil {
		return err
	}

	return s.tokens.RevokeFamily(ctx, record.FamilyID)
}

func (s *Service) issueTokens(
	ctx context.Context,
	user *models.User,
	familyID, userAgent, ip string,
) (*TokenPair, error) {
	access, err := utils.GenerateAccessToken(
		user.ID.Hex(), string(user.Role), s.jwtSecret,
	)
	if err != nil {
		return nil, err
	}

	refresh, err := utils.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	if familyID == "" {
		familyID = bson.NewObjectID().Hex()
	}

	return &TokenPair{access, refresh}, s.tokens.Create(
		ctx,
		&models.RefreshToken{
			UserID:    user.ID,
			FamilyID:  familyID,
			TokenHash: utils.HashToken(refresh),
			UserAgent: userAgent,
			IPAddress: ip,
			ExpiresAt: time.Now().UTC().Add(utils.RefreshTokenTTL),
		},
	)
}

func (s *Service) findToken(
	ctx context.Context, token string,
) (*models.RefreshToken, error) {
	record, err := s.tokens.FindByHash(ctx, utils.HashToken(token))
	if errors.Is(err, repo.ErrTokenNotFound) {
		return nil, ErrInvalidToken
	}
	if err != nil {
		return nil, err
	}
	if time.Now().After(record.ExpiresAt) {
		return nil, ErrInvalidToken
	}

	return record, nil
}