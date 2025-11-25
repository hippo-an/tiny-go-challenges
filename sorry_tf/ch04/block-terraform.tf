terraform {
    required_version = ">= 1.11.0"

    required_providers {
        aws = {
            source = "hashicorp/aws"
            version = "~> 5.0"
        }
    }

    backend "s3" {
        baucket = "terraform-state-backend"
        key = "state/terraform.tfstate"
        region = "ap-northeast-2"
    }
}