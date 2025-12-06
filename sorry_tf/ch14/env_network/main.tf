# remote backend
locals {
  topology = yamldecode(file("./topology.yaml"))

  tf_vpc_env_list = distinct(flatten([
    for k, v in local.topology : [v.requester.tf_env, v.accepter.tf_env]
  ]))
}

data "terraform_remote_state" "vpc" {
  for_each = toset(local.tf_vpc_env_list)
  backend  = "s3"
  config = {
    bucket  = "terraform-tfstate"
    key     = "${each.key}.tfstate"
    region  = "ap-northeast-2"
    profile = "terraform"
  }
}

locals {
  vpc_ids = {
    for k, v in data.terraform_remote_state.vpc : k => v.outputs.vpc_id
  }
}

# vpc peering
locals {
  env_tags = {
    tf_env = "env_network"
  }
}

module "seoul_to_virginia_peering" {
  source = "../../ch11/modules/vpc_peering"
  for_each = {
    for k, v in local.topology : k => v
    if v.requester.tf_env == "seoul" && v.accepter.tf_env == "virginia"
  }

  providers = {
    aws.requester = aws.seoul
    aws.accepter  = aws.virginia
  }

  name             = each.key
  requester_vpc_id = local.vpc_ids[each.value.requester.tf_env][each.value.requester.vpc]
  accepter_vpc_id  = local.vpc_ids[each.value.accepter.tf_env][each.value.accepter.vpc]
  tags             = local.env_tags
}

module "seoul_to_seoul_peering" {
  source = "../../ch11/modules/vpc_peering"
  for_each = {
    for k, v in local.topology : k => v
    if v.requester.tf_env == "seoul" && v.accepter.tf_env == "seoul"
  }

  providers = {
    aws.requester = aws.seoul
    aws.accepter  = aws.seoul
  }

  name             = each.key
  requester_vpc_id = local.vpc_ids[each.value.requester.tf_env][each.value.requester.vpc]
  accepter_vpc_id  = local.vpc_ids[each.value.accepter.tf_env][each.value.accepter.vpc]
  tags             = local.env_tags
}
