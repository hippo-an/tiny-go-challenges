# modules/nested-module/main.tf
resource "aws_instance" "example" {
  ami           = "ami-123123123123"
  instance_type = "t3.micro"
  tags = {
    Name = var.instance_name
  }
}

resource "aws_s3_bucket" "example" {
  bucket = var.bucket_name
}

resource "aws_s3_bucket_acl" "example" {
  bucket = aws_s3_bucket.example.id
  acl    = "private"
}

# modules/nested-module/variables.tf
variable "instance_name" {
  type = string
}

variable "bucket_name" {
  type = string
}

# modules/nested-module/outputs.tf
output "instance_id" {
  value = aws_instance.example.id
}

output "bucket_arn" {
  value = aws_s3_bucket.example.arn
}
