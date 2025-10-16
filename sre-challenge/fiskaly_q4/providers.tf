terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 6.0"
    }
  }
}

# Override default region that I have locally set to "eu-central-1"
provider "aws" {
  region = var.aws_region
}