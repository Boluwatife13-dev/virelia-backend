package utils

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"
)

const (
	OTPTTL         = 10 * time.Minute
	MaxOTPAttempts = 5
)

func GenerateOTP() (string, error) {
	max := big.NewInt(1000000)

	number, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%06d", number.Int64()), nil
}

func HashOTP(code, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(code))

	return hex.EncodeToString(mac.Sum(nil))
}

func CheckOTP(code, hash, secret string) bool {
	expected := HashOTP(code, secret)

	return hmac.Equal(
		[]byte(expected),
		[]byte(hash),
	)
}