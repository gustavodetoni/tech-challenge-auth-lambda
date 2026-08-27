output "function_name" {
  description = "Nome da Lambda de autenticacao CPF/CNPJ."
  value       = aws_lambda_function.auth_cpf.function_name
}

output "function_arn" {
  description = "ARN da Lambda de autenticacao CPF/CNPJ."
  value       = aws_lambda_function.auth_cpf.arn
}

output "function_invoke_arn" {
  description = "Invoke ARN usado pelo API Gateway."
  value       = aws_lambda_function.auth_cpf.invoke_arn
}

output "security_group_id" {
  description = "Security group usado pela Lambda."
  value       = data.terraform_remote_state.k8s.outputs.auth_lambda_security_group_id
}

