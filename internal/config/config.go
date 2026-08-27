package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	DatabaseURL      string
	JWTSecret        string
	JWTIssuer        string
	JWTAudience      string
	JWTExpiryMinutes int
}

func Load() (Config, error) {
	cfg := Config{
		DatabaseURL:      strings.TrimSpace(os.Getenv("DATABASE_URL")),
		JWTSecret:        strings.TrimSpace(os.Getenv("JWT_SECRET")),
		JWTIssuer:        envOrDefault("JWT_ISSUER", "tech-challenge-auth-lambda"),
		JWTAudience:      strings.TrimSpace(os.Getenv("JWT_AUDIENCE")),
		JWTExpiryMinutes: envInt("JWT_EXPIRY_MINUTES", 60),
	}

	if cfg.DatabaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required")
	}
	if cfg.JWTSecret == "" {
		return Config{}, fmt.Errorf("JWT_SECRET is required")
	}
	if cfg.JWTExpiryMinutes < 1 {
		return Config{}, fmt.Errorf("JWT_EXPIRY_MINUTES must be >= 1")
	}

	return cfg, nil
}

func envOrDefault(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func envInt(key string, fallback int) int {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}
