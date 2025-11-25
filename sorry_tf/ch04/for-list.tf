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

> [for k, v in aws_vpc.this : v.id]
[
    "vpc-123123123123",
    "vpc-456456456456",
    "vpc-789789789789"
]