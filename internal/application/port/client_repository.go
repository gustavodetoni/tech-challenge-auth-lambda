package port

import (
	"context"

	"github.com/soat-architecture/tech-challenge-auth-lambda/internal/domain"
)

type ClientRepository interface {
	FindByDocument(ctx context.Context, documentNumber string) (*domain.Client, error)
}
