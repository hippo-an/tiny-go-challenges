# access key
provider "aws" {
    region = "ap-northeast-2"
    access_key = "xxxxx"
    secret_key = "yyyyy"
}

# profile
provider "aws" {
    region = "ap-northeast-2"
    profile = "terraform"
}

# assume role
provider "aws" {
    region = "ap-northeast-2"
    profile = "terraform"

    assume_role = {
        role_arn = "arn:aws:iam::${account_id}:role/AssumeRole"
    }
}

# AWS 기본 프로바이더
provider "aws" {
    region = "ap-northeast-2"
    profile = "terraform-a"
}

# AWS 추가 프로바이더
provider "aws" {
    region = "ap-northeast-1"
    profile = "terraform-b"
    alias = "terraform-b"
}

# 기본 프로바이더 사용
resource "aws_instance" "this" {
    provider = "aws" # 생략가능
    instance_type = "t3.medium"
}

# 추가 프로바이더 사용
resource "aws_instance" "this" {
    provider = aws.terraform-b
    instance_type = "t3.medium"
}

# 기본 프로바이더 사용
module "my_module" {
    source = "./modules/my_module"
    providers = {
        aws = aws
    }
}

# 추가 프로바이더 사용
module "my_module" {
    source = "./modules/my_module"
    providers = {
        aws = aws.terraform-b
    }
}
