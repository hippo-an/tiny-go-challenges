terraform {
  backend "s3" {
    bucket  = "terraform-tfstate"
    key     = "virginia.tfstate"
    region  = "ap-northeast-2"
    encrypt = true
    profile = "terraform"
  }
}
