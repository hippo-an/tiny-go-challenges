# tags
locals {
  vpc_name = var.vpc_name
  vpc_id   = var.vpc_id

  module_tag = merge(
    var.tags,
    {
      tf_module = "ec2"
    }
  )
}

locals {
  valid_ebs_type = ["standard", "gp2", "gp3", "io1", "io2", "sc1", "st1"]
}

# ec2 인스턴스
resource "aws_instance" "this" {
  for_each = var.ec2_set

  # required
  ami                    = each.value.ami
  instance_type          = each.value.instance_type
  subnet_id              = var.subnet_id_map[each.value.subnet][each.value.az]
  vpc_security_group_ids = [for sg_name in each.value.security_group : var.sg_id_map[sg_name]]

  # optional
  iam_instance_profile = each.value.ec2_role
  key_name             = each.value.ec2_key
  private_ip           = each.value.private_ip

  # root volume
  root_block_device {
    volume_type           = each.value.root_volume.type
    volume_size           = each.value.root_volume.size
    delete_on_termination = true

    tags = merge(
      local.ec2_tags[each.key],
      {
        Name = "${each.value.full_name}-root"
      }
    )
  }

  # tags
  tags = local.ec2_tags[each.key]

  # lifecycle
  lifecycle {
    precondition {
      condition     = contains(["develop", "staging", "rc", "production"], each.value.env)
      error_message = "[${local.vpc_name}] env must be one of [develop, staging, rc, production]"
    }

    precondition {
      condition     = contains(local.valid_ebs_type, each.value.root_volume.type)
      error_message = "[${local.vpc_name} VPC/${each.key}] root_volume.type must be one of ${join(", ", local.valid_ebs_type)}"
    }
  }
}

# ec2_set override
locals {
  ec2_set = {
    for k, v in var.ec2_set : k => merge(v, {
      full_name = "${var.vpc_name}-${split("-", v.subnet)[0]}-${k}"
    })
  }

  ec2_tags = {
    for k, v in local.ec2_set : k => merge(
      local.module_tag,
      {
        Name    = v.full_name
        EC2     = v.full_name
        Env     = v.env
        Team    = v.team
        Service = v.service
        OS      = upper(v.os_type)
      }
    )
  }
}

# EIP 리소스 블록 정의
resource "aws_eip" "this" {
  for_each = {
    for k, v in local.ec2_set : k => v
    if split("-", v.subnet)[0] == "pub"
  }

  domain = "vpc"

  instance                  = aws_instance.this[each.key].id
  associate_with_private_ip = aws_instance.this[each.key].private_ip

  tags = local.ec2_tags[each.key]
}

# additional ebs volume info prepare
locals {
  ec2_volume_set = [
    for ec2_name, ec2_attribute in local.ec2_set : {
      for volume in ec2_attribute.additional_volume : "${ec2_name}_${volume.device}" => merge(
        { ec2_name = ec2_name }, volume
      )
    }
  ]

  merged_ec2_volume_set = module.merge_ec2_volume_set.output
}

module "merge_ec2_volume_set" {
  source = "../../../ch09/merge_map_module"
  input  = local.ec2_volume_set
}

# additional ebs volume
locals {
  valid_iops_type = ["gp3", "io1", "io2"]
}

resource "aws_ebs_volume" "this" {
  for_each          = local.merged_ec2_volume_set
  availability_zone = aws_instance.this[each.value.ec2_name].availability_zone
  size              = each.value.size
  type              = each.value.type
  iops              = contains(local.valid_iops_type, each.value.type) ? each.value.iops : null

  tags = merge(
    local.ec2_tags[each.value.ec2_name],
    {
      Name = "${each.value.full_name}-${each.value.device}"
    }
  )

  # lifecycle
  lifecycle {

    precondition {
      condition     = contains(local.valid_ebs_type, each.value.type)
      error_message = "[${local.vpc_name} VPC/${each.key} EC2:${each.value.device} EBS] additional_volume.type must be one of ${join(", ", local.valid_ebs_type)}"
    }

    precondition {
      condition     = !(!contains(local.valid_iops_type, each.value.type) && each.value.iops != null)
      error_message = "[${local.vpc_name} VPC/${each.key} EC2:${each.value.device} EBS] additional_volume.iops must be null if additional_volume.type is not one of ${join(", ", local.valid_iops_type)}"
    }

  }
}

resource "aws_volume_attachment" "this" {
  for_each    = local.merged_ec2_volume_set
  device_name = each.value.device
  volume_id   = aws_ebs_volume.this[each.key].id
  instance_id = aws_instance.this[each.value.ec2_name].id
}
