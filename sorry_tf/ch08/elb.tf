variable "lb_info" {
  default = {
    name = "test-lb"
    type = "application"
  }
}

variable "listener_info" {
  default = {
    port        = 443
    protocol    = "HTTPS"
    rules       = {}
    alpn_policy = "None"
  }

}

variable "lb_type" {
  default = "application"
}

locals {
  is_alb = var.lb_info.type == "application"
  is_nlb = var.lb_info.type == "network"

  listener_rules = var.listener_info.rules
}

resource "aws_lb_listener" "this" {
  alpn_policy = var.protocol == "TLS" ? var.listener_info.alpn_policy : null
}

resource "aws_lb_listener_rule" "this" {
  for_each = {
    for rule in var.listener_info.rules : rule.name => rule
    if local.is_alb
  }
}

check "nlb_listener_have_rule" {
  assert {
    condition     = !(local.is_nlb && length(var.listener_info.rules) > 0)
    error_message = "NLB listener should not have rules"
  }
}

check "alb_listener_have_aplpn_policy" {
  assert {
    condition     = !(local.is_alb && var.listener_info.alpn_policy != "None")
    error_message = "ALB listener should not have ALPN policy"
  }
}
