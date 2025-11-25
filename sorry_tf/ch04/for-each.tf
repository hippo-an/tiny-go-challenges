resource "aws_instance" "this" {
    for_each = {
        windows = {
            ami = "ami-0123456789"
            type = "t3.micro"
        }

        linux = {
            ami = "ami-9876543210"
            type = "r5.2xlarge"
        }
    }

    ami = each.value.ami
    instanct_type = each.value.type

    subnet_id = "subnet-123456"
    vpc_security_group_ids = ["sg-123456"]

    tags = {
        Name = each.key
    }

}