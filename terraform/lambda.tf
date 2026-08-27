locals {
  function_name = "${var.project_name}-${var.environment}-auth-cpf"
  database_url  = "postgres://${var.db_user}:${var.db_password}@${data.terraform_remote_state.database.outputs.database_host}:${data.terraform_remote_state.database.outputs.database_port}/${data.terraform_remote_state.database.outputs.database_name}?sslmode=require"
}

resource "aws_iam_role" "lambda" {
  name = "${local.function_name}-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      }
    ]
  })

  tags = merge(var.tags, {
    Environment = var.environment
  })
}

resource "aws_iam_role_policy_attachment" "basic_execution" {
  role       = aws_iam_role.lambda.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_iam_role_policy_attachment" "vpc_execution" {
  role       = aws_iam_role.lambda.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

resource "aws_lambda_function" "auth_cpf" {
  function_name = local.function_name
  role          = aws_iam_role.lambda.arn

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

  depends_on = [
    aws_iam_role_policy_attachment.basic_execution,
    aws_iam_role_policy_attachment.vpc_execution,
  ]
}

