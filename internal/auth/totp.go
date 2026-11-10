package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"net/url"
	"time"
)

var b32 = base32.StdEncoding.WithPadding(base32.NoPadding)

func newTOTP(email string) (secret, link string, err error) {
	b := make([]byte, 20)

	if _, err = rand.Read(b); err != nil {
		return "", "", err
	}

	secret = b32.EncodeToString(b)

	link = fmt.Sprintf(
		"otpauth://totp/%s?secret=%s&issuer=Virelia",
		url.PathEscape("Virelia:"+email), secret,
	)

	return secret, link, nil
}

func checkTOTP(secret, code string) bool {
	key, err := b32.DecodeString(secret)
	if err != nil {
		return false
	}

	now := time.Now().Unix() / 30

	for step := now - 1; step <= now+1; step++ {
		msg := binary.BigEndian.AppendUint64(nil, uint64(step))

		mac := hmac.New(sha1.New, key)
		mac.Write(msg)
		sum := mac.Sum(nil)

		i := sum[len(sum)-1] & 0x0f
		n := binary.BigEndian.Uint32(sum[i:i+4]) & 0x7fffffff

		if fmt.Sprintf("%06d", n%1000000) == code {
			return true
		}
	}

	return false
}
