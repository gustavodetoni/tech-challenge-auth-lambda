terraform {
  required_version = ">= 1.6.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = var.aws_region
}

data "aws_caller_identity" "current" {}

locals {
  aws_academy_service_role_arn = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:role/${var.aws_academy_service_role_name}"
  lambda_role_arn              = var.lambda_role_arn != "" ? var.lambda_role_arn : local.aws_academy_service_role_arn
}

data "terraform_remote_state" "k8s" {
  backend = "s3"

  config = {
    bucket = var.k8s_state_bucket
    key    = var.k8s_state_key
    region = var.k8s_state_region
  }
}

data "terraform_remote_state" "database" {
  backend = "s3"

  config = {
    bucket = var.database_state_bucket
    key    = var.database_state_key
    region = var.database_state_region
  }
}
