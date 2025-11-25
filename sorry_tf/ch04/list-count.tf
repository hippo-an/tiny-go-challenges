locals {
    ami_list = ["ami-123456", "ami-456789", "ami-987654"]
}

resource "aws_instance" "this" {
    count = length(local.ami_list)  # length 로 list 길이 입력
    ami = local.ami_list[count.index]  # list 의 값을 index 로 가지고 온다.
    instance_type = "t3.micro"

    private_ip = "10.0.0.${count.index + 1}"
    
    tags = {
        Name = "EC2-${count.index + 1}"
    }
}