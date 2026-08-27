package port

import (
	"time"

	"github.com/soat-architecture/tech-challenge-auth-lambda/internal/domain"
)

type ClientTokenClaims struct {
	ClientID       string
	DocumentType   domain.DocumentType
	DocumentNumber string
	Scopes         []string
}

type TokenIssuer interface {
	NewClientToken(clientEntity domain.Client) (string, time.Time, error)
}

type TokenValidator interface {
	ParseAndValidate(tokenString string) (*ClientTokenClaims, error)
}
