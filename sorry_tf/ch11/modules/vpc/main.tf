#common tags ===========================================
locals {
  vpc_name = var.name

  module_tag = merge(
    var.tags,
    {
      tf_module = "vpc"
      Env       = var.attribute.env
      Team      = var.attribute.team
      VPC       = "${local.vpc_name}-vpc"
    }
  )
}

# vpc ===========================================
locals {
  vpc_cidr = var.attribute.cidr
  vpc_id   = aws_vpc.this.id
}

resource "aws_vpc" "this" {
  cidr_block           = local.vpc_cidr
  enable_dns_support   = true
  enable_dns_hostnames = true
  tags = merge(
    local.module_tag,
    {
      Name = "${local.vpc_name}-vpc"
    }
  )
}

# public subnets ===========================================
locals {
  subnets = var.attribute.subnets
  enable_igw = anytrue(
    [for k, v in local.subnets : split("-", k)[0] == "pub"]
  )
}

resource "aws_internet_gateway" "this" {
  vpc_id = local.vpc_id
  count  = local.enable_igw ? 1 : 0

  tags = merge(
    local.module_tag,
    {
      Name = "${local.vpc_name}-igw"
    }
  )
}


resource "aws_route_table" "public" {
  count  = local.enable_igw ? 1 : 0
  vpc_id = local.vpc_id

  tags = merge(
    local.module_tag,
    {
      Name = "${local.vpc_name}-rt-pub"
    }
  )
}

resource "aws_route" "public_igw" {
  count                  = local.enable_igw ? 1 : 0
  destination_cidr_block = "0.0.0.0/0"
  route_table_id         = aws_route_table.public[count.index].id
  gateway_id             = aws_internet_gateway.this[count.index].id
}

# private route table ===========================================
locals {
  subnet_azs = var.attribute.subnet_azs
}

resource "aws_route_table" "private" {
  for_each = toset(local.subnet_azs)

  vpc_id = local.vpc_id

  tags = merge(
    local.module_tag,
    {
      Name = "${local.vpc_name}-rt-pri-${each.value}"
    }
  )
}

# subnets ===========================================
locals {
  subnet_newbits = var.attribute.subnet_newbits
  subnet_azs     = var.attribute.subnet_azs

  subnets_data = flatten([
    for name, indices in local.subnets : [
      for idx in indices : {
        name      = name
        az        = local.subnet_azs[index(indices, idx)]
        cidr      = cidrsubnet(local.vpc_cidr, local.subnet_newbits, idx)
        is_public = split("-", name)[0] == "pub"
      }
    ]
  ])
}

locals {
  subnets_map = {
    for s in local.subnets_data : "${replace(s.name, "-", "")}_${s.az}" => s
  }
}

module "current" {
  source = "../../../ch09/meta_module"
}

locals {
  region_name   = module.current.region_name
  available_azs = module.current.az_names
}

resource "aws_subnet" "this" {
  for_each = local.subnets_map

  vpc_id                  = local.vpc_id
  cidr_block              = each.value.cidr
  availability_zone       = "${local.region_name}${each.value.az}"
  map_public_ip_on_launch = each.value.is_public

  tags = merge(
    local.module_tag,
    {
      Name = "${local.vpc_name}-subnet-${each.value.name}-${each.value.az}"
    }
  )

  lifecycle {
    precondition {
      condition     = contains(local.available_azs, "${local.region_name}${each.value.az}")
      error_message = "${upper(each.value.az)} zone 은 현재 리전 (${local.region_name}) 에서 유효하지 않다. 사용 가능 영역: [${join(", ", [for az in local.available_azs : trimprefix(az, local.region_name)])}]"
    }
  }

  lifecycle {
    precondition {
      condition     = contains(["pub", "pri"], split("-", each.value.name)[0])
      error_message = "[${local.vpc_name}] ${each.value.name} 이라는 서브넷 이름은 유효하지 않습니다. subnets 이름들은 모두 [pub-, pri-]로 시작해야 합니다."
    }
  }
}

# route table association ===========================================
locals {
  public_rt = try(aws_route_table.public[0].id, "")
  private_rts = {
    for k, v in aws_route_table.private : k => v.id
  }
}

resource "aws_route_table_association" "this" {
  for_each = local.subnets_map

  subnet_id      = aws_subnet.this[each.key].id
  route_table_id = each.value.is_public ? local.public_rt : local.private_rts[each.value.az]
}


# mapping ==========================================
locals {
  subnet_ids = {
    for k, v in aws_subnet.this : k => v.id
  }

  subnet_ids_with_az = {
    for k, v in aws_subnet.this : k => {
      for az in slice(local.subnet_azs, 0, length(v)) : az => aws_subnet.this["${replace(k, "-", "_")}_${az}"].id
    }
  }

}
