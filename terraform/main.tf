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

