variable "input" {
  description = "list(map()) or map(map())"
  type        = list(map(string))
}

locals {
  keys   = flatten([for item in var.input : keys(item)])
  values = flatten([for item in var.input : values(item)])

  output = zipmap(local.keys, local.values)
}

output "output" {
  value = local.output
}

