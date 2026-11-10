package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var ErrInvalidToken = errors.New("invalid or expired token")

const (
	loginIssuer     = "virelia"
	challengeIssuer = "virelia-2fa"
)

type Claims struct {
	Role       string `json:"role"`
	FacilityID string `json:"facility_id,omitempty"`
	TOTPSetup  string `json:"totp_setup,omitempty"`
	jwt.RegisteredClaims
}

func signToken(
	c Claims, userID, issuer, secret string, ttl time.Duration,
) (string, error) {
	c.Subject = userID
	c.Issuer = issuer
	c.ExpiresAt = jwt.NewNumericDate(time.Now().Add(ttl))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	return token.SignedString([]byte(secret))
}

func parseToken(token, issuer, secret string) (*Claims, error) {
	claims := &Claims{}

	key := func(*jwt.Token) (any, error) { return []byte(secret), nil }

	_, err := jwt.ParseWithClaims(token, claims, key,
		jwt.WithValidMethods([]string{"HS256"}),
		jwt.WithIssuer(issuer),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

func ParseAccessToken(token, secret string) (*Claims, error) {
	return parseToken(token, loginIssuer, secret)
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))

	return hex.EncodeToString(sum[:])
}

func (s *Service) issueTokens(
	ctx context.Context, userID,
	family,
	agent,
	ip string,
) (*TokenPair, error) {
	role, err := s.repo.RoleOf(ctx, userID)
	if err != nil {
		return nil, ErrInvalidToken
	}

	claims := Claims{Role: role.Name, FacilityID: role.FacilityID}

	access, err := signToken(
		claims, userID, loginIssuer, s.secret, 15*time.Minute,
	)
	if err != nil {
		return nil, err
	}

	b := make([]byte, 32)
	rand.Read(b)
	refresh := hex.EncodeToString(b)

	t := RefreshToken{UserID: userID, FamilyID: family}

	err = s.repo.SaveRefreshToken(ctx, t, hashToken(refresh), agent, ip)

	return &TokenPair{access, refresh}, err
}
