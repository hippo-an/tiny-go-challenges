locals {
  info_files = "${path.root}/../info_files"

  vpc_set = toset([
    for vpcfile in fileset(local.info_files, "*/vpc.yaml") : dirname(vpcfile)
  ])

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

module "sg" {
  for_each = local.vpc_set
  source   = "../modules/sg"

  vpc_name = each.key
  vpc_id   = module.vpc[each.key].vpc_id

  sg_set = {
    for sgfile in fileset(local.info_files, "${each.key}/sg/*.csv") :
    trimsuffix(basename(sgfile), ".csv") => csvdecode(file("${local.info_files}/${sgfile}"))
  }

  tags = local.env_tags
}

module "ec2" {
  for_each = local.vpc_set
  source   = "../modules/ec2"

  vpc_name = each.key
  vpc_id   = module.vpc[each.key].vpc_id

  subnet_id_map = module.vpc[each.key].subnet_ids_with_az
  sg_id_map     = module.sg[each.key].sg_id

  ec2_set = {
    for ec2file in fileset(local.info_files, "${each.key}/ec2/*.yaml") :
    trimsuffix(basename(ec2file), ".yaml") => yamldecode(file("${local.info_files}/${ec2file}"))
  }

  tags = local.env_tags
}
