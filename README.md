# Tech Challenge Auth Lambda

Function serverless responsavel por autenticar clientes por CPF/CNPJ e emitir JWTs para consumo das APIs protegidas da oficina.

## Repositorios Da Entrega

- Aplicacao principal: https://github.com/gustavodetoni/tech-challenge-project
- Lambda Auth CPF/CNPJ: https://github.com/gustavodetoni/tech-challenge-auth-lambda
- Infra Kubernetes: https://github.com/gustavodetoni/tech-challenge-infra-k8s
- Infra Database: https://github.com/gustavodetoni/tech-challenge-infra-database

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

O deploy e executado manualmente pelo GitHub Actions para facilitar a demonstracao e o destroy no AWS Academy.
O gatilho automatico por `push` esta comentado no workflow e deve ser habilitado apenas quando as branches de homologacao/producao estiverem configuradas.

Fluxo previsto:

```text
pull_request -> lint/test
Run workflow -> action=apply, environment=homolog
Run workflow -> action=destroy, environment=homolog
```

Inputs do workflow manual:

```text
action       apply ou destroy
environment  homolog ou prod
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

Secrets necessarios:

```text
AWS_ACCESS_KEY_ID
AWS_SECRET_ACCESS_KEY
AWS_SESSION_TOKEN
TF_STATE_BUCKET
TF_VAR_DB_PASSWORD
JWT_SECRET
```

O workflow cria o bucket de state automaticamente caso ele ainda nao exista.
O workflow imprime `function_name` e `function_invoke_arn` no resumo do GitHub Actions.
Esses valores devem ser usados no segundo `apply` do repositorio `tech-challenge-infra-k8s`.

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

- Repositorio: https://github.com/gustavodetoni/tech-challenge-auth-lambda
- Swagger da API principal: https://github.com/gustavodetoni/tech-challenge-project/blob/main/docs/swagger.yaml
- Postman da API principal: https://github.com/gustavodetoni/tech-challenge-project/blob/main/docs/collections/tech-challenge.postman_collection.json
- Deploy homologacao: sera atualizado apos o primeiro deploy cloud.
- Deploy producao: sera atualizado apos o primeiro deploy cloud.
