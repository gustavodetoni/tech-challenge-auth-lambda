package db

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	repository "github.com/soat-architecture/tech-challenge-auth-lambda/internal/application/port"
	"github.com/soat-architecture/tech-challenge-auth-lambda/internal/domain"
)

type ClientRepository struct {
	pool *pgxpool.Pool
}

func NewClientRepository(pool *pgxpool.Pool) *ClientRepository {
	return &ClientRepository{pool: pool}
}

func (r *ClientRepository) FindByDocument(ctx context.Context, documentNumber string) (*domain.Client, error) {
	const query = `
		SELECT
			id::text,
			document_type::text,
			document_number,
			name,
			CASE WHEN deleted_at IS NULL THEN 'ACTIVE' ELSE 'INACTIVE' END AS status
		FROM clients
		WHERE document_number = $1
		ORDER BY deleted_at NULLS FIRST
		LIMIT 1;
	`

	var out domain.Client
	err := r.pool.QueryRow(ctx, query, documentNumber).Scan(
		&out.ID,
		&out.DocumentType,
		&out.DocumentNumber,
		&out.Name,
		&out.Status,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}

	return &out, nil
}
