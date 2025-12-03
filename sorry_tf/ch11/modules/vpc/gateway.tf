# NAT gateway ===========================================
locals {
  nat     = var.attribute.naaat
  nat_azs = slice(local.subnet_azs, 0, local.nat.per_az ? try(length(local.subnets[local.nat.subnet]), 0) : 1)
  nat_set = local.nat.create ? toset(local.nat_azs) : toset([])
}

## nat eip
resource "aws_eip" "this" {
  for_each = local.nat_set
  tags = merge(
    local.module_tag,
    {
      Name = "${local.vpc_name}-nat-${each.key}"
    }
  )
}

## nat gateway
resource "aws_nat_gateway" "this" {
  for_each      = local.nat_set
  allocation_id = aws_eip.this[each.key].id
  subnet_id     = local.subnet_ids_with_az[local.nat.subnet][each.key]

  tags = merge(
    local.module_tag,
    {
      Name = "${local.vpc_name}-nat-${each.key}"
    }
  )

  lifecycle {
    precondition {
      condition     = split("-", local.nat.subnet)[0] == "pub"
      error_message = "[${local.vpc_name}] nat.subnet 으로는 퍼블릭 서브넷만 지정 가능합니다."
    }
  }
}

resource "aws_route" "private_nat" {
  for_each = local.nat.create ? local.private_rts : {}

  route_table_id         = each.value
  destination_cidr_block = "0.0.0.0/0"
  nat_gateway_id         = aws_nat_gateway.this[element(local.nat_azs, index(local.subnet_azs, each.key) % length(local.nat_azs))].id
}

