data "aws_vpc" "requester" {
  provider = aws.requester
  id       = var.requester_vpc_id
}

data "aws_vpc" "accepter" {
  provider = aws.accepter
  id       = var.accepter_vpc_id
}

locals {
  requester_vpc = data.aws_vpc.requester
  accepter_vpc  = data.aws_vpc.accepter

  module_tag = merge(
    var.tags,
    {
      Name         = var.name
      tf_moudle    = "vpc_peering"
      Request_VPC  = lookup(local.requester_vpc.tags, "Name", "네임태그없음")
      Accepter_VPC = lookup(local.accepter_vpc.tags, "Name", "네임태그없음")
    }
  )
}

# vpc peering 연결 생성 및 자동 수락
module "accepter" {
  source = "../../../ch09/meta_module"
  providers = {
    aws = aws.accepter
  }
}

locals {
  accepter_account_id = module.accepter.account_id
  accepter_region     = module.accepter.region_name
}


# route table 
data "aws_route_tables" "requester" {
  provider = aws.requester
  vpc_id   = var.requester_vpc_id
}

data "aws_route_tables" "accepter" {
  provider = aws.accepter
  vpc_id   = var.accepter_vpc_id
}

locals {
  requester_rtbs = data.aws_route_tables.requester.ids
  accepter_rtbs  = data.aws_route_tables.accepter.ids
}
