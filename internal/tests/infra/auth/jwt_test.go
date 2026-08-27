package auth_test

import (
	"testing"
	"time"

	"github.com/soat-architecture/tech-challenge-auth-lambda/internal/domain"
	infraAuth "github.com/soat-architecture/tech-challenge-auth-lambda/internal/infra/auth"
)

func TestManagerIssuesAndValidatesClientToken(t *testing.T) {
	manager := infraAuth.NewManager("secret", "issuer", "audience", time.Hour)

	token, _, err := manager.NewClientToken(domain.Client{
		ID:             "client-1",
		DocumentType:   domain.DocumentTypeCPF,
		DocumentNumber: "52998224725",
		Status:         domain.StatusActive,
	})
	if err != nil {
		t.Fatalf("NewClientToken() error = %v", err)
	}

	claims, err := manager.ParseAndValidate(token)
	if err != nil {
		t.Fatalf("ParseAndValidate() error = %v", err)
	}
	if claims.ClientID != "client-1" {
		t.Fatalf("ClientID = %q", claims.ClientID)
	}
	if claims.DocumentNumber != "52998224725" {
		t.Fatalf("DocumentNumber = %q", claims.DocumentNumber)
	}
}
