package interfaces

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"

	appAuth "github.com/soat-architecture/tech-challenge-auth-lambda/internal/application/auth"
	repository "github.com/soat-architecture/tech-challenge-auth-lambda/internal/application/port"
)

type AuthService interface {
	AuthenticateDocument(ctx context.Context, input appAuth.AuthenticateDocumentInput) (*appAuth.AuthenticateDocumentOutput, error)
	ValidateBearerToken(token string) (*repository.ClientTokenClaims, error)
}

type Handler struct {
	auth AuthService
}

func NewHandler(auth AuthService) *Handler {
	return &Handler{auth: auth}
}

func (h *Handler) Handle(ctx context.Context, payload json.RawMessage) (any, error) {
	var probe map[string]any
	if err := json.Unmarshal(payload, &probe); err != nil {
		return nil, err
	}
	if _, ok := probe["routeArn"]; ok {
		return h.handleAuthorizer(ctx, payload)
	}
	return h.handleAPI(ctx, payload)
}

func (h *Handler) handleAPI(ctx context.Context, payload json.RawMessage) (events.APIGatewayV2HTTPResponse, error) {
	var request events.APIGatewayV2HTTPRequest
	if err := json.Unmarshal(payload, &request); err != nil {
		return jsonResponse(http.StatusBadRequest, errorResponse{Error: "invalid request"}), nil
	}

	if request.RequestContext.HTTP.Method != http.MethodPost {
		return jsonResponse(http.StatusMethodNotAllowed, errorResponse{Error: "method not allowed"}), nil
	}

	body, err := decodeBody(request.Body, request.IsBase64Encoded)
	if err != nil {
		return jsonResponse(http.StatusBadRequest, errorResponse{Error: "invalid request body"}), nil
	}

	var input documentAuthRequest
	if err := json.Unmarshal([]byte(body), &input); err != nil {
		return jsonResponse(http.StatusBadRequest, errorResponse{Error: "invalid request body"}), nil
	}

	out, err := h.auth.AuthenticateDocument(ctx, appAuth.AuthenticateDocumentInput{Document: input.Document})
	if err != nil {
		return jsonResponse(statusFromError(err), errorResponse{Error: publicMessage(err)}), nil
	}

	return jsonResponse(http.StatusOK, tokenResponse{
		AccessToken: out.AccessToken,
		TokenType:   out.TokenType,
		ExpiresAt:   out.ExpiresAt,
		Client: clientResponse{
			ID:       out.Client.ID,
			Document: out.Client.Document,
			Type:     string(out.Client.Type),
			Status:   string(out.Client.Status),
		},
	}), nil
}

func (h *Handler) handleAuthorizer(ctx context.Context, payload json.RawMessage) (authorizerResponse, error) {
	var request authorizerRequest
	if err := json.Unmarshal(payload, &request); err != nil {
		return authorizerResponse{IsAuthorized: false}, nil
	}

	token := bearerToken(headerValue(request.Headers, "authorization"))
	claims, err := h.auth.ValidateBearerToken(token)
	if err != nil {
		return authorizerResponse{IsAuthorized: false}, nil
	}

	claimsJSON, _ := json.Marshal(claims)
	return authorizerResponse{
		IsAuthorized: true,
		Context: map[string]any{
			"client_claims": string(claimsJSON),
		},
	}, nil
}

func decodeBody(body string, encoded bool) (string, error) {
	if !encoded {
		return body, nil
	}
	decoded, err := base64.StdEncoding.DecodeString(body)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}

func jsonResponse(status int, body any) events.APIGatewayV2HTTPResponse {
	data, _ := json.Marshal(body)
	return events.APIGatewayV2HTTPResponse{
		StatusCode: status,
		Headers: map[string]string{
			"content-type": "application/json",
		},
		Body: string(data),
	}
}

func statusFromError(err error) int {
	switch {
	case errors.Is(err, appAuth.ErrInvalidInput):
		return http.StatusBadRequest
	case errors.Is(err, appAuth.ErrUnauthorized):
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}

func publicMessage(err error) string {
	switch {
	case errors.Is(err, appAuth.ErrInvalidInput):
		return "invalid document"
	case errors.Is(err, appAuth.ErrUnauthorized):
		return "unauthorized"
	default:
		return "internal server error"
	}
}

func bearerToken(value string) string {
	parts := strings.Fields(value)
	if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
		return parts[1]
	}
	return ""
}

func headerValue(headers map[string]string, key string) string {
	for currentKey, value := range headers {
		if strings.EqualFold(currentKey, key) {
			return value
		}
	}
	return ""
}

type documentAuthRequest struct {
	Document string `json:"document"`
}

type tokenResponse struct {
	AccessToken string         `json:"access_token"`
	TokenType   string         `json:"token_type"`
	ExpiresAt   string         `json:"expires_at"`
	Client      clientResponse `json:"client"`
}

type clientResponse struct {
	ID       string `json:"id"`
	Document string `json:"document"`
	Type     string `json:"type"`
	Status   string `json:"status"`
}

type errorResponse struct {
	Error string `json:"error"`
}

type authorizerRequest struct {
	RouteArn string            `json:"routeArn"`
	Headers  map[string]string `json:"headers"`
}

type authorizerResponse struct {
	IsAuthorized bool           `json:"isAuthorized"`
	Context      map[string]any `json:"context,omitempty"`
}
