variable "project_name" {
  description = "Nome base usado nos recursos criados na AWS."
  type        = string
  default     = "tech-challenge"
}

variable "aws_region" {
  description = "Regiao AWS onde os recursos serao criados."
  type        = string
  default     = "us-east-1"
}

variable "environment" {
  description = "Ambiente de deploy."
  type        = string
  default     = "homolog"
}

variable "lambda_package_path" {
  description = "Caminho do zip gerado pelo build da Lambda."
  type        = string
  default     = "../dist/function.zip"
}

variable "db_user" {
  description = "Usuario do PostgreSQL usado pela Lambda."
  type        = string
  default     = "tech_challenge"
}

variable "db_password" {
  description = "Senha do PostgreSQL usada pela Lambda."
  type        = string
  sensitive   = true
}

variable "jwt_secret" {
  description = "Secret HS256 usado para assinar os JWTs de cliente."
  type        = string
  sensitive   = true
}

variable "jwt_expiry_minutes" {
  description = "Tempo de validade do JWT em minutos."
  type        = number
  default     = 60
}

variable "k8s_state_bucket" {
  description = "Bucket S3 do state do repositorio tech-challenge-infra-k8s."
  type        = string
}

variable "k8s_state_key" {
  description = "Key do state do repositorio tech-challenge-infra-k8s."
  type        = string
  default     = "tech-challenge/k8s/terraform.tfstate"
}

variable "k8s_state_region" {
  description = "Regiao do bucket S3 do state do repositorio tech-challenge-infra-k8s."
  type        = string
  default     = "us-east-1"
}

variable "database_state_bucket" {
  description = "Bucket S3 do state do repositorio tech-challenge-infra-database."
  type        = string
}

variable "database_state_key" {
  description = "Key do state do repositorio tech-challenge-infra-database."
  type        = string
  default     = "tech-challenge/database/terraform.tfstate"
}

variable "database_state_region" {
  description = "Regiao do bucket S3 do state do repositorio tech-challenge-infra-database."
  type        = string
  default     = "us-east-1"
}

variable "tags" {
  description = "Tags aplicadas aos recursos."
  type        = map(string)

  default = {
    Project = "tech-challenge"
    Managed = "terraform"
  }
}

