data "aws-ami" "this" {
  owners      = ["amazon"]
  most_recent = true
  filter {
    name   = "image-id"
    values = ["ami-0c55b159cbfafe1f0"]
  }
}

resource "aws_instance" "this" {
  ami           = data.aws_ami.this.id
  instance_type = "t2.micro"

  lifecycle {
    precondition {
      condition     = data.aws_ami.this.architecture == "x86_64"
      error_message = "ami architecture must be x86_64"
    }
    postcondition {
      condition     = self.public_dns != ""
      error_message = "instance must have a public dns"
    }
  }
}
