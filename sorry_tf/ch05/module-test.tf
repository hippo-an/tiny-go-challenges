module "instance" {
  source = "./ec2-asg"
  # github.com/hippo-an/tf-ec2-asg 와 같이 github 형식 가능

  minimum_count = 3
  maximum_count = 50
  desired_count = 4

  instance_type = "t2.micro"
}

module "instance2" {
  source  = "./ec2-asg"
  version = "5.8.1"
  # version = "~> 5.8.1" # 비관적 제약조건

  minimum_count = 3
  maximum_count = 50
  desired_count = 4

  instance_type = "t2.micro"
}

provider "aws" {
  region = "ap-northeast-2"
  alias  = "seoul"
}

provider "aws" {
  region = "ap-northeast-1"
  alias  = "tokyo"
}

module "example" {
  source = "./ec2-asg"
  providers = {
    aws = aws.seoul
    aws = aws.tokyo
  }
}
