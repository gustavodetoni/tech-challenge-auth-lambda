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
- Datadog/CloudWatch para logs e rastreabilidade operacional

## Governanca Dos Repositorios

Os quatro repositorios da entrega possuem branch protection ativa na branch `main`, exigindo Pull Request para merge, execucao das validacoes de CI e impedindo commits diretos como fluxo oficial de desenvolvimento.

Branches de homologacao e producao sao atendidas por GitHub Actions. O deploy e automatizado pela esteira: apos o disparo definido no workflow, a pipeline empacota o binario Go, aplica o Terraform e publica a Lambda na AWS sem comandos manuais nos servidores.

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

O deploy e automatizado pelo GitHub Actions para AWS usando Terraform.
A esteira executa testes, empacota a Lambda, inicializa o backend remoto S3, gera o plano Terraform e aplica a infraestrutura da function serverless.

Antes do `terraform init`, o workflow faz bootstrap do backend S3. Esse passo cria o bucket de state quando necessario e cria um state vazio valido quando o objeto `tech-challenge/lambda/<ambiente>.tfstate` tiver sido removido, evitando a falha `HeadObject 403` comum em contas AWS Academy sem permissao de listagem completa.

Fluxo previsto:

```text
pull_request      -> lint/test
main/homolog/prod -> package + terraform apply
destroy           -> terraform destroy controlado
```

Inputs do workflow:

```text
action       apply ou destroy
environment  homolog ou prod
```

O Terraform da Lambda fica na pasta `terraform/` e consome os outputs dos repositorios:

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
Em AWS Academy, a Lambda reutiliza por padrao a role pre-criada `LabRole`, evitando criacao/anexo de IAM policies pela esteira. Se outro ambiente exigir uma role diferente, informe `TF_VAR_lambda_role_arn` ou ajuste `lambda_role_arn` no Terraform.
O workflow imprime `function_name` e `function_invoke_arn` no resumo do GitHub Actions.
Esses valores devem ser usados no segundo `apply` do repositorio `tech-challenge-infra-k8s`.

## Observabilidade

A Lambda participa da rastreabilidade do fluxo de autenticacao por CPF/CNPJ.
O token emitido carrega os dados do cliente e e usado pelas rotas protegidas da API principal, que registra `traceId`/`correlation_id` nos logs JSON.

Na AWS, a execucao da Lambda gera logs operacionais no CloudWatch, enquanto o cluster EKS usa Datadog Agent para coletar metricas, logs estruturados, latencia, erros e traces da aplicacao principal. Dessa forma, o fluxo completo pode ser demonstrado desde `POST /auth/cpf` ate o consumo das APIs protegidas.

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
- Endpoint homologacao: https://hwq42fgalh.execute-api.us-east-1.amazonaws.com/auth/cpf
- Deploy producao: mesmo fluxo automatizado de deploy, usando `environment=prod`.
