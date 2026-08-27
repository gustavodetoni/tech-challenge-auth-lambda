package auth_test

import (
	"context"
	"errors"
	"testing"
	"time"

	appAuth "github.com/soat-architecture/tech-challenge-auth-lambda/internal/application/auth"
	repository "github.com/soat-architecture/tech-challenge-auth-lambda/internal/application/port"
	"github.com/soat-architecture/tech-challenge-auth-lambda/internal/domain"
)

func TestAuthenticateDocumentIssuesClientToken(t *testing.T) {
	repo := fakeClientRepository{client: &domain.Client{
		ID:             "client-1",
		DocumentType:   domain.DocumentTypeCPF,
		DocumentNumber: "52998224725",
		Status:         domain.StatusActive,
	}}
	issuer := fakeTokenManager{token: "signed", expiresAt: time.Date(2026, 8, 26, 22, 0, 0, 0, time.UTC)}
	service := appAuth.NewClientAuthUseCase(repo, issuer, issuer)

	out, err := service.AuthenticateDocument(context.Background(), appAuth.AuthenticateDocumentInput{Document: "529.982.247-25"})
	if err != nil {
		t.Fatalf("AuthenticateDocument() error = %v", err)
	}

	if out.AccessToken != "signed" {
		t.Fatalf("AccessToken = %q", out.AccessToken)
	}
	if out.Client.ID != "client-1" {
		t.Fatalf("Client.ID = %q", out.Client.ID)
	}
}

func TestAuthenticateDocumentRejectsUnknownClient(t *testing.T) {
	repo := fakeClientRepository{err: repository.ErrNotFound}
	issuer := fakeTokenManager{}
	service := appAuth.NewClientAuthUseCase(repo, issuer, issuer)

	_, err := service.AuthenticateDocument(context.Background(), appAuth.AuthenticateDocumentInput{Document: "52998224725"})
	if !errors.Is(err, appAuth.ErrUnauthorized) {
		t.Fatalf("error = %v, want ErrUnauthorized", err)
	}
}

func TestAuthenticateDocumentRejectsInvalidCPF(t *testing.T) {
	issuer := fakeTokenManager{}
	service := appAuth.NewClientAuthUseCase(fakeClientRepository{}, issuer, issuer)

	_, err := service.AuthenticateDocument(context.Background(), appAuth.AuthenticateDocumentInput{Document: "11111111111"})
	if !errors.Is(err, appAuth.ErrInvalidInput) {
		t.Fatalf("error = %v, want ErrInvalidInput", err)
	}
}

type fakeClientRepository struct {
	client *domain.Client
	err    error
}

func (f fakeClientRepository) FindByDocument(context.Context, string) (*domain.Client, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.client, nil
}

type fakeTokenManager struct {
	token     string
	expiresAt time.Time
	claims    *repository.ClientTokenClaims
}

func (f fakeTokenManager) NewClientToken(domain.Client) (string, time.Time, error) {
	return f.token, f.expiresAt, nil
}

func (f fakeTokenManager) ParseAndValidate(string) (*repository.ClientTokenClaims, error) {
	return f.claims, nil
}
