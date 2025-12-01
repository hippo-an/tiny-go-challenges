variable "eks_version" {
  type    = string
  default = "1.30"
}

variable "eks_addon_name" {
  type    = string
  default = "vpc-cni"
}

variable "eks_addon_version" {
  type    = string
  default = "1.18.1-eksbuild.3"
}

data "aws_eks_addon_version" "this" {
  kubernetes_version = var.eks_version
  addon_name         = var.eks_addon_name
}


locals {
  addon_default_version = data.aws_eks_addon_version.this.version
}

check "addon_version" {
    assert {
        condition = var.eks_addon_version == local.addon_default_version
        error_message = "recommended version is ${local.addon_default_version}"
    }
}
