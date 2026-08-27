# Tech Challenge Auth Lambda

Function serverless responsavel por autenticar clientes por CPF/CNPJ e emitir JWTs para consumo das APIs protegidas da oficina.

## Proposito

Esta Lambda atende o requisito de autenticacao serverless do desafio. Ela valida o documento informado, consulta a base de clientes no banco gerenciado e devolve um token JWT com escopos de cliente.

## Tecnologias

- Go
- AWS Lambda
- AWS API Gateway
- PostgreSQL gerenciado
- JWT
- GitHub Actions

## Contrato Inicial

```http
POST /auth/cpf
Content-Type: application/json
```

```json
{
  "document": "12345678909"
}
```

Resposta esperada:

```json
{
  "access_token": "jwt",
  "token_type": "Bearer",
  "expires_at": "2026-08-26T22:00:00Z",
  "client": {
    "id": "uuid",
    "document": "12345678909",
    "status": "ACTIVE"
  }
}
```

## Execucao Local

```bash
go test ./...
go run ./cmd/auth
```

Variaveis esperadas:

```env
DATABASE_URL=postgres://user:pass@host:5432/tech_challenge?sslmode=require
JWT_SECRET=change-me
JWT_ISSUER=tech-challenge-auth-lambda
JWT_EXPIRY_MINUTES=60
```

## Deploy

O deploy sera executado por GitHub Actions nas branches de homologacao e producao.

Fluxo previsto:

```text
pull_request -> lint/test
homolog      -> deploy homolog
main         -> deploy producao
```

## Arquitetura

```text
API Gateway
  -> Lambda Auth CPF
      -> RDS PostgreSQL
      -> JWT assinado
```

## Links

- Swagger/Postman da API principal: pendente
- Deploy homologacao: pendente
- Deploy producao: pendente

