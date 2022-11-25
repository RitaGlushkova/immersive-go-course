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

provider "aws" {
  profile = "cyfplus"
  region = "us-east-1"
  default_tags {
    tags = {
      Name = "RitaGlushkova/immersive-go-course"
      owner= "RitaGlushkova"
      project = "docker-cloud"
    }
  }
}

provider "aws" {
  alias = "global_region"
  profile = "cyfplus"
  region = "eu-west-2"
  default_tags {
    tags = {
      Name = "RitaGlushkova/immersive-go-course"
      owner= "RitaGlushkova"
      project = "docker-cloud"
    }
  }
}