resource "aws_vpc" "this" {
    for_each = toset(["a", "b", "c"])
    cidr_block = "10.0.0.0/16"
}

> aws_vpc.this["a"]
{
    "arn" = "arn:aws:ec2:ap-northeast-2:xxxx"
    "cidr_block" = "10.0.0.0/16"
    "id" = "vpc-123123123123"
}

> {for k, v in aws_vpc.this : k => v.id}
{
    "a" = "vpc-123123123123",
    "b" = "vpc-456456456456",
    "c" = "vpc-789789789789"
}

> {for k, v in aws_vpc.this : k => {id = v.id, cidr = v.cidr_block}}
{
    "a" = {
        id = "vpc-123123123123"
        cidr = "10.0.0.0/16"
    },
    "b" = {
        id = "vpc-456456456456",
        cidr = "10.0.0.0/16"
    }
    "c" = {
        id = "vpc-789789789789"
        cidr = "10.0.0.0/16"
    } 
}