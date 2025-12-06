# Requester side
resource "aws_vpc_peering_connection" "this" {
  provider = aws.requester

  vpc_id = local.requester_vpc.id

  peer_vpc_id   = local.accepter_vpc.id
  peer_owner_id = local.accepter_account_id
  peer_region   = local.accepter_region

  auto_accept = false

  tags = local.module_tag
}

# Accepter side
resource "aws_vpc_peering_connection_accepter" "this" {
  provider = aws.accepter

  vpc_peering_connection_id = aws_vpc_peering_connection.this.id
  auto_accept               = true

  tags = local.module_tag
}

# VPC Peering Options : DNS Resolution
locals {
  peering_id = aws_vpc_peering_connection_accepter.this.id
}

resource "aws_vpc_peering_connection_options" "requester" {
  provider = aws.requester

  vpc_peering_connection_id = local.peering_id

  requester {
    allow_dns_resolution_from_remote_vpc = true
  }
}

resource "aws_vpc_peering_connection_options" "accepter" {
  provider = aws.accepter

  vpc_peering_connection_id = local.peering_id

  accepter {
    allow_dns_resolution_from_remote_vpc = true
  }
}

# VPC Route table
resource "aws_route" "requester_to_accepter" {
  for_each = toset(local.requester_rtbs)
  provider = aws.requester

  route_table_id            = each.key
  destination_cidr_block    = local.accepter_vpc.cidr_block
  vpc_peering_connection_id = local.peering_id
}

resource "aws_route" "accepter_to_requester" {
  for_each = toset(local.accepter_rtbs)
  provider = aws.accepter

  route_table_id            = each.key
  destination_cidr_block    = local.requester_vpc.cidr_block
  vpc_peering_connection_id = local.peering_id
}
