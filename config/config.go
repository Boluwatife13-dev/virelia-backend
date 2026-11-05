package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port         string
	MongoURI     string
	DatabaseName string
	JWTSecret    string
	BrevoAPIKey  string
	EmailFrom    string
	FrontendURLs string
}

func LoadConfig() Config {
	_ = godotenv.Load()

	jwtSecret := os.Getenv("JWT_SECRET")
	if len(jwtSecret) < 32 {
		log.Fatal("JWT_SECRET must be at least 32 characters")
	}

	return Config{
		Port:         getEnv("PORT", "8081"),
		MongoURI:     getEnv("MONGODB_URI", "mongodb://localhost:27017"),
		DatabaseName: getEnv("DATABASE_NAME", "virelia"),
		JWTSecret:    jwtSecret,
		BrevoAPIKey:  os.Getenv("BREVO_API_KEY"),
		EmailFrom:    os.Getenv("EMAIL_FROM"),
		FrontendURLs: getEnv("FRONTEND_URLS", "http://localhost:3000"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}