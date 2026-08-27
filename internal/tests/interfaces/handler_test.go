package interfaces_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/aws/aws-lambda-go/events"

	appAuth "github.com/soat-architecture/tech-challenge-auth-lambda/internal/application/auth"
	repository "github.com/soat-architecture/tech-challenge-auth-lambda/internal/application/port"
	"github.com/soat-architecture/tech-challenge-auth-lambda/internal/domain"
	lambdaAdapter "github.com/soat-architecture/tech-challenge-auth-lambda/internal/interfaces"
)

func TestHandleAPIAuthenticatesDocument(t *testing.T) {
	handler := lambdaAdapter.NewHandler(fakeAuthService{
		output: &appAuth.AuthenticateDocumentOutput{
			AccessToken: "token",
			TokenType:   "Bearer",
			ExpiresAt:   "2026-08-26T22:00:00Z",
			Client: appAuth.ClientOutput{
				ID:       "client-1",
				Document: "52998224725",
				Type:     domain.DocumentTypeCPF,
				Status:   domain.StatusActive,
			},
		},
	})

	event := events.APIGatewayV2HTTPRequest{
		Body: `{"document":"529.982.247-25"}`,
		RequestContext: events.APIGatewayV2HTTPRequestContext{
			HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
				Method: "POST",
			},
		},
	}
	payload, _ := json.Marshal(event)

	response, err := handler.Handle(context.Background(), payload)
	if err != nil {
		t.Fatalf("Handle() error = %v", err)
	}

	apiResponse := response.(events.APIGatewayV2HTTPResponse)
	if apiResponse.StatusCode != 200 {
		t.Fatalf("StatusCode = %d", apiResponse.StatusCode)
	}
}

func TestHandleAuthorizerRejectsMissingToken(t *testing.T) {
	handler := lambdaAdapter.NewHandler(fakeAuthService{})
	payload := []byte(`{"routeArn":"arn","headers":{}}`)

	response, err := handler.Handle(context.Background(), payload)
	if err != nil {
		t.Fatalf("Handle() error = %v", err)
	}

	data, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	var body struct {
		IsAuthorized bool `json:"isAuthorized"`
	}
	if err := json.Unmarshal(data, &body); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if body.IsAuthorized {
		t.Fatal("expected unauthorized response")
	}
}

type fakeAuthService struct {
	output *appAuth.AuthenticateDocumentOutput
}

func (f fakeAuthService) AuthenticateDocument(context.Context, appAuth.AuthenticateDocumentInput) (*appAuth.AuthenticateDocumentOutput, error) {
	return f.output, nil
}

func (f fakeAuthService) ValidateBearerToken(string) (*repository.ClientTokenClaims, error) {
	return nil, appAuth.ErrUnauthorized
}
