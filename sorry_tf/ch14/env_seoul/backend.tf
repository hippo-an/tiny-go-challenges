terraform {
  backend "s3" {
    bucket  = "terraform-tfstate"
    key     = "seoul.tfstate"
    region  = "ap-northeast-2"
    encrypt = true
    profile = "terraform"
  }
}
