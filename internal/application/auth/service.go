package auth

import (
	"context"
	"errors"
	"fmt"

	repository "github.com/soat-architecture/tech-challenge-auth-lambda/internal/application/port"
	"github.com/soat-architecture/tech-challenge-auth-lambda/internal/domain"
	"github.com/soat-architecture/tech-challenge-auth-lambda/pkg/br/document"
)

var (
	ErrInvalidInput = errors.New("invalid input")
	ErrUnauthorized = errors.New("unauthorized")
)

type ClientAuthUseCase struct {
	clients repository.ClientRepository
	issuer  repository.TokenIssuer
	tokens  repository.TokenValidator
}

func NewClientAuthUseCase(clients repository.ClientRepository, issuer repository.TokenIssuer, tokens repository.TokenValidator) *ClientAuthUseCase {
	return &ClientAuthUseCase{clients: clients, issuer: issuer, tokens: tokens}
}

type AuthenticateDocumentInput struct {
	Document string
}

type AuthenticateDocumentOutput struct {
	AccessToken string
	TokenType   string
	ExpiresAt   string
	Client      ClientOutput
}

type ClientOutput struct {
	ID       string
	Document string
	Type     domain.DocumentType
	Status   domain.Status
}

func (s *ClientAuthUseCase) AuthenticateDocument(ctx context.Context, input AuthenticateDocumentInput) (*AuthenticateDocumentOutput, error) {
	doc := document.Normalize(input.Document)
	docType, err := documentType(doc)
	if err != nil {
		return nil, err
	}

	clientEntity, err := s.clients.FindByDocument(ctx, doc)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrUnauthorized
		}
		return nil, err
	}

	if clientEntity.DocumentType != docType || !clientEntity.IsActive() {
		return nil, ErrUnauthorized
	}

	token, expiresAt, err := s.issuer.NewClientToken(*clientEntity)
	if err != nil {
		return nil, err
	}

	return &AuthenticateDocumentOutput{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresAt:   expiresAt.UTC().Format("2006-01-02T15:04:05Z07:00"),
		Client: ClientOutput{
			ID:       clientEntity.ID,
			Document: clientEntity.DocumentNumber,
			Type:     clientEntity.DocumentType,
			Status:   clientEntity.Status,
		},
	}, nil
}

func (s *ClientAuthUseCase) ValidateBearerToken(token string) (*repository.ClientTokenClaims, error) {
	if token == "" {
		return nil, ErrUnauthorized
	}
	claims, err := s.tokens.ParseAndValidate(token)
	if err != nil {
		return nil, ErrUnauthorized
	}
	return claims, nil
}

func documentType(doc string) (domain.DocumentType, error) {
	switch {
	case document.IsValidCPF(doc):
		return domain.DocumentTypeCPF, nil
	case document.IsValidCNPJ(doc):
		return domain.DocumentTypeCNPJ, nil
	default:
		return "", fmt.Errorf("%w: invalid document", ErrInvalidInput)
	}
}
