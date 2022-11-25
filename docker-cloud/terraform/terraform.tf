terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.0"
    }
  }
  backend "s3" {
    bucket  = "rita-terraform-state"
    key     = "state.tfstate"
    region  = "eu-west-2"
    profile = "cyfplus"
  }

  
}

locals {
  tags = {
    project = "docker-cloud"
    terraform   = "true"
    owner       = "rita"
  }
}

provider "aws" {
  region  = "eu-west-2"
  profile = "cyfplus"

  default_tags {
    tags = local.tags
  }
}