# meta_module
module "current" {
  source = "./meta_module"
}

locals {
  account_id        = module.current.account_id
  region            = module.current.region_name
  account_alias     = module.current.account_alias
  region_code       = module.current.region_code
  availability_zone = module.current.az_name
}


# provider_validation
module "check_cross" {
  source = "./provider_validation_module"
  providers = {
    aws.a = aws.requester
    aws.b = aws.accepter
  }
}


locals {
  is_cross_account = module.check_cross.is_cross_account
  is_cross_region  = module.check_cross.is_cross_region

  need_accepter = local.is_cross_account || local.is_cross_region
}

resource "aws_vpc_peering_connection" "this" {
  peer_owner_id = local.is_cross_account ? local.accepter_account : null
  peer_region   = local.is_cross_region ? local.accepter_region : null

  auto_accept = local.need_accepter ? false : true

  dynamic "requester" {
    for_each = local.need_accepter ? toset([]) : toset(["1"])
    content {
      allow_remote_vpc_dns_resolution = true

    }
  }

  dynamic "accepter" {
    for_each = local.need_accepter ? toset([]) : toset(["1"])
    content {
      allow_remote_vpc_dns_resolution = true
    }
  }

  # ...
}

# merge_map_module
locals {
  vpc_list = [
    {
      vpc1 = "1234",
      vpc2 = "5678"
    },
    {
      vpc3 = "9876"
    }
  ]

  vpc_map = merge(local.vpc_list)
}

module "merge_vpc_list" {
  source = "./merge_map_module"
  input  = local.vpc_list
}

module "merge_vpc_list_from_map" {
  source = "./merge_map_module"
  input  = local.vpc_map
}

output "output" {
  value = module.merge_vpc_list.output
}

output "output_from_map" {
  value = module.merge_vpc_list_from_map.output
}
