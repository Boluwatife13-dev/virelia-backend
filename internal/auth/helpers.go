package auth

import (
	"crypto/rand"
	"strings"
)

func newNetworkID() string {
	const chars = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"

	b := make([]byte, 8)
	rand.Read(b)

	for i := range b {
		b[i] = chars[int(b[i])%len(chars)]
	}

	return "VRL-" + string(b)
}

func normalize(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func optional(value string) *string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}

	return &value
}
