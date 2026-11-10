package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
)

var ErrBadKey = errors.New("ENCRYPTION_KEY must be 64 hex characters")

var ErrBadData = errors.New("could not decrypt")

type Box struct {
	gcm cipher.AEAD
}

// New makes a Box from a 64-character hex key.
func New(hexKey string) (*Box, error) {
	key, err := hex.DecodeString(hexKey)

	if err != nil || len(key) != 32 {
		return nil, ErrBadKey
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return &Box{gcm}, nil
}

func (b *Box) Seal(text string) (string, error) {
	nonce := make([]byte, b.gcm.NonceSize())

	_, err := rand.Read(nonce)
	if err != nil {
		return "", err
	}

	sealed := b.gcm.Seal(nonce, nonce, []byte(text), nil)

	return base64.StdEncoding.EncodeToString(sealed), nil
}

func (b *Box) Open(sealed string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(sealed)

	size := b.gcm.NonceSize()

	if err != nil || len(data) < size {
		return "", ErrBadData
	}

	nonce, body := data[:size], data[size:]

	text, err := b.gcm.Open(nil, nonce, body, nil)
	if err != nil {
		return "", ErrBadData
	}

	return string(text), nil
}
