locals {
  info_files = "${path.root}/../info_files"

  vpc_set = toset([
    for vpcfile in fileset(local.info_files, "*/vpc.yaml") : dirname(vpcfile)
  ])
}

locals {
  env_tags = {
    tf_env = "ch11/env"
  }
}

module "vpc" {
  for_each = local.vpc_set
  source   = "../modules/vpc"

  name      = each.key
  attribute = yamldecode(file("${local.info_files}/${each.key}/vpc.yaml"))

  tags = local.env_tags
}
