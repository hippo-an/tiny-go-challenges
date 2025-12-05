# security group =========================================
locals {
  vpc_name = var.vpc_name
  vpc_id   = var.vpc_id
  vpc_tags = data.aws_vpc.this.tags

  module_tag = merge(
    var.tags,
    local.vpc_tags,
    {
      tf_module = "sg"
    }
  )
}


data "aws_vpc" "this" {
  id = local.vpc_id
}


locals {
  tf_desc = "Managed By Terraform"
}

resource "aws_security_group" "this" {
  for_each    = var.sg_set
  name        = "${local.vpc_name}-sg-${each.key}"
  description = local.tf_desc
  vpc_id      = local.vpc_id

  tags = local.module_tag
}

# egress rule =========================================
resource "aws_security_group_rule" "this" {
  for_each          = var.sg_set
  security_group_id = aws_security_group.this[each.key].id
  type              = "egress"
  from_port         = 0
  to_port           = 0
  protocol          = "-1"
  cidr_blocks       = ["0.0.0.0/0"]
  description       = local.tf_desc
}

# utility inbound_rule_set =========================================
locals {
  inbound_rule_set = [
    for sg, rules in var.sg_set : {
      for r in rules : "${sg}-${r.protocol}-${r.from_port}-${r.to_port}-${r.source}" => merge(r, { sg = sg })
    }
  ]
}

module "merge_inbound_rule_set" {
  source = "../../../ch09/merge_map_module"
  input  = local.inbound_rule_set
}

locals {
  merged_inbound_rule_set = module.merge_inbound_rule_set.output
}


# ingress rule =========================================# 
resource "aws_security_group_rule" "this" {
  for_each          = local.merged_inbound_rule_set
  security_group_id = aws_security_group.this[each.value.sg].id
  type              = "ingress"
  from_port         = each.value.from_port
  to_port           = each.value.to_port
  protocol          = each.value.protocol

  self                     = each.value.source == "self" ? true : null
  cidr_blocks              = length(regexall("[a-z]", each.value.source)) == 0 ? [each.value.source] : null
  source_security_group_id = startswith(each.value.source, "sg-") ? each.value.source : null
  prefix_list_ids          = startswith(each.value.source, "pl-") ? [each.value.source] : null
  description              = each.value.desc == "" ? "tf/${each.value.source}" : "tf/${each.value.desc}"
}

