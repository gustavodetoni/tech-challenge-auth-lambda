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
- Terraform
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
make tests
make package
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

O Terraform da Lambda fica em `terraform/` e consome os outputs dos repositorios:

- `tech-challenge-infra-k8s`: VPC, subnets privadas e security group da Lambda.
- `tech-challenge-infra-database`: endpoint, porta e nome do RDS PostgreSQL.

Exemplo:

```bash
cp terraform/terraform.tfvars.example terraform/terraform.tfvars
make package
terraform -chdir=terraform init
terraform -chdir=terraform plan
terraform -chdir=terraform apply
```

## Arquitetura

```text
API Gateway
  -> Lambda Auth CPF
      -> RDS PostgreSQL
      -> JWT assinado
```

Internamente, o codigo segue a mesma separacao da API principal:

```text
cmd                       Composition root da Lambda
internal/domain           Entidades centrais
internal/application      Caso de uso e portas
internal/interfaces       Adapter AWS Lambda/API Gateway
internal/infra            JWT e PostgreSQL
pkg/br                    Validacao reutilizavel de CPF/CNPJ
terraform                 Infraestrutura da Lambda
```

## Links

- Swagger/Postman da API principal: pendente
- Deploy homologacao: pendente
- Deploy producao: pendente
