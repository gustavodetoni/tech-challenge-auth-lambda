package config_test

import (
	"testing"

	"github.com/soat-architecture/tech-challenge-auth-lambda/internal/config"
)

func TestLoadRequiresDatabaseURL(t *testing.T) {
	t.Setenv("DATABASE_URL", "")
	t.Setenv("JWT_SECRET", "secret")

	_, err := config.Load()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestLoadWithDefaults(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://example")
	t.Setenv("JWT_SECRET", "secret")
	t.Setenv("JWT_ISSUER", "")
	t.Setenv("JWT_EXPIRY_MINUTES", "")

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.JWTIssuer != "tech-challenge-auth-lambda" {
		t.Fatalf("JWTIssuer = %q", cfg.JWTIssuer)
	}
	if cfg.JWTExpiryMinutes != 60 {
		t.Fatalf("JWTExpiryMinutes = %d", cfg.JWTExpiryMinutes)
	}
}
