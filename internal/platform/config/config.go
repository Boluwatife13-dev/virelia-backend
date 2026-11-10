package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port          string
	DatabaseURL   string
	JWTSecret     string
	FrontendURLs  []string
	BrevoAPIKey   string
	EmailFrom     string
	EncryptionKey string
	WebAuthnRPID  string
}

func Load() Config {
	_ = godotenv.Load()

	cfg := Config{
		Port:          getEnv("PORT", "8081"),
		DatabaseURL:   os.Getenv("DATABASE_URL"),
		JWTSecret:     os.Getenv("JWT_SECRET"),
		BrevoAPIKey:   os.Getenv("BREVO_API_KEY"),
		EmailFrom:     os.Getenv("EMAIL_FROM"),
		EncryptionKey: os.Getenv("ENCRYPTION_KEY"),
		WebAuthnRPID:  getEnv("WEBAUTHN_RP_ID", "localhost"),
		FrontendURLs: strings.Split(
			getEnv("FRONTEND_URLS", "http://localhost:3000"),
			",",
		),
	}

	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	if len(cfg.JWTSecret) < 32 {
		log.Fatal("JWT_SECRET must be at least 32 characters")
	}

	if len(cfg.EncryptionKey) != 64 {
		log.Fatal("ENCRYPTION_KEY must be 64 hex characters")
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
