locals {
  function_name = "${var.project_name}-${var.environment}-auth-cpf"
  database_url  = "postgres://${var.db_user}:${var.db_password}@${data.terraform_remote_state.database.outputs.database_host}:${data.terraform_remote_state.database.outputs.database_port}/${data.terraform_remote_state.database.outputs.database_name}?sslmode=require"
}

resource "aws_lambda_function" "auth_cpf" {
  function_name = local.function_name
  role          = local.lambda_role_arn

  filename         = var.lambda_package_path
  source_code_hash = filebase64sha256(var.lambda_package_path)

  runtime = "provided.al2023"
  handler = "bootstrap"
  timeout = 10

  environment {
    variables = {
      DATABASE_URL       = local.database_url
      JWT_SECRET         = var.jwt_secret
      JWT_ISSUER         = "tech-challenge-auth-lambda"
      JWT_EXPIRY_MINUTES = tostring(var.jwt_expiry_minutes)
    }
  }

  vpc_config {
    subnet_ids         = data.terraform_remote_state.k8s.outputs.private_subnet_ids
    security_group_ids = [data.terraform_remote_state.k8s.outputs.auth_lambda_security_group_id]
  }

  tags = merge(var.tags, {
    Environment = var.environment
  })
}
