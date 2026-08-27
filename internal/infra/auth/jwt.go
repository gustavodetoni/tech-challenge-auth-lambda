package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	repository "github.com/soat-architecture/tech-challenge-auth-lambda/internal/application/port"
	"github.com/soat-architecture/tech-challenge-auth-lambda/internal/domain"
)

type Claims struct {
	DocumentType   string   `json:"document_type"`
	DocumentNumber string   `json:"document"`
	TokenType      string   `json:"type"`
	Scopes         []string `json:"scope"`
	jwt.RegisteredClaims
}

type Manager struct {
	secret   []byte
	issuer   string
	audience string
	ttl      time.Duration
}

func NewManager(secret, issuer, audience string, ttl time.Duration) *Manager {
	return &Manager{
		secret:   []byte(secret),
		issuer:   issuer,
		audience: audience,
		ttl:      ttl,
	}
}

func (m *Manager) NewClientToken(clientEntity domain.Client) (string, time.Time, error) {
	if len(m.secret) == 0 {
		return "", time.Time{}, fmt.Errorf("JWT secret not configured")
	}
	if clientEntity.ID == "" {
		return "", time.Time{}, fmt.Errorf("client id is required")
	}

	now := time.Now().UTC()
	expiresAt := now.Add(m.ttl)
	scopes := []string{"client:service-orders:read", "client:budget:decide"}

	claims := Claims{
		DocumentType:   string(clientEntity.DocumentType),
		DocumentNumber: clientEntity.DocumentNumber,
		TokenType:      "CLIENT",
		Scopes:         scopes,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   clientEntity.ID,
			Issuer:    m.issuer,
			Audience:  jwt.ClaimStrings(nonEmptyString(m.audience)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(m.secret)
	if err != nil {
		return "", time.Time{}, err
	}

	return signed, expiresAt, nil
}

func (m *Manager) ParseAndValidate(tokenString string) (*repository.ClientTokenClaims, error) {
	if len(m.secret) == 0 {
		return nil, fmt.Errorf("JWT secret not configured")
	}
	if tokenString == "" {
		return nil, fmt.Errorf("missing token")
	}

	claims := &Claims{}
	parserOpts := []jwt.ParserOption{
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithIssuedAt(),
	}
	if m.issuer != "" {
		parserOpts = append(parserOpts, jwt.WithIssuer(m.issuer))
	}
	if m.audience != "" {
		parserOpts = append(parserOpts, jwt.WithAudience(m.audience))
	}

	parser := jwt.NewParser(parserOpts...)
	token, err := parser.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return m.secret, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid || claims.TokenType != "CLIENT" {
		return nil, fmt.Errorf("invalid token")
	}

	return &repository.ClientTokenClaims{
		ClientID:       claims.Subject,
		DocumentType:   domain.DocumentType(claims.DocumentType),
		DocumentNumber: claims.DocumentNumber,
		Scopes:         claims.Scopes,
	}, nil
}

func nonEmptyString(value string) []string {
	if value == "" {
		return nil
	}
	return []string{value}
}
