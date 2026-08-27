package main

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/lambda"

	appAuth "github.com/soat-architecture/tech-challenge-auth-lambda/internal/application/auth"
	"github.com/soat-architecture/tech-challenge-auth-lambda/internal/config"
	infraAuth "github.com/soat-architecture/tech-challenge-auth-lambda/internal/infra/auth"
	"github.com/soat-architecture/tech-challenge-auth-lambda/internal/infra/db"
	lambdaAdapter "github.com/soat-architecture/tech-challenge-auth-lambda/internal/interfaces"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := db.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	clientRepo := db.NewClientRepository(pool)
	tokenManager := infraAuth.NewManager(cfg.JWTSecret, cfg.JWTIssuer, cfg.JWTAudience, time.Duration(cfg.JWTExpiryMinutes)*time.Minute)
	authUseCase := appAuth.NewClientAuthUseCase(clientRepo, tokenManager, tokenManager)
	handler := lambdaAdapter.NewHandler(authUseCase)

	lambda.Start(handler.Handle)
}
