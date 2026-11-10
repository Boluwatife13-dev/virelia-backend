package auth

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"time"
)

const (
	otpTTL      = 10 * time.Minute
	maxAttempts = 5
)

var ErrInvalidCode = errors.New(
	"invalid or expired code",
)

func newCode() (string, error) {
	max := big.NewInt(1000000)

	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}

	code := fmt.Sprintf("%06d", n.Int64())

	return code, nil
}

func hashCode(code, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))

	mac.Write([]byte(code))

	hash := mac.Sum(nil)

	return hex.EncodeToString(hash)
}

func (r *Repository) SaveOTP(
	ctx context.Context,
	userID string,
	purpose string,
	hash string,
) error {
	query := `
		INSERT INTO otp_codes (
			user_id,
			purpose,
			code_hash,
			expires_at
		)
		VALUES ($1, $2, $3, $4)

		ON CONFLICT (user_id, purpose)
		DO UPDATE SET
			code_hash  = $3,
			expires_at = $4,
			attempts   = 0
	`

	expiresAt := time.Now().Add(otpTTL)

	_, err := r.db.Exec(
		ctx,
		query,
		userID,
		purpose,
		hash,
		expiresAt,
	)

	return err
}

func (r *Repository) CheckOTP(
	ctx context.Context,
	userID string,
	purpose string,
	hash string,
) error {
	query := `
		UPDATE otp_codes
		SET attempts = attempts + 1

		WHERE user_id  = $1
		  AND purpose  = $2
		  AND attempts < $3

		RETURNING id, code_hash, expires_at
	`

	var id string
	var savedHash string
	var expiresAt time.Time

	err := r.db.QueryRow(
		ctx,
		query,
		userID,
		purpose,
		maxAttempts,
	).Scan(&id, &savedHash, &expiresAt)

	if errors.Is(dbError(err), ErrNotFound) {
		return ErrInvalidCode
	}

	if err != nil {
		return err
	}

	matches := hmac.Equal(
		[]byte(savedHash),
		[]byte(hash),
	)

	expired := time.Now().After(expiresAt)

	if !matches || expired {
		return ErrInvalidCode
	}

	_, err = r.db.Exec(
		ctx,
		"DELETE FROM otp_codes WHERE id = $1",
		id,
	)

	return err
}
