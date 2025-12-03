variable "name" {
  type        = string
  description = "구분되는 VPC 이름"
}

variable "attribute" {
  description = "VPC 속성 정의"
  type = object({
    cidr           = string
    env            = optional(string, "develop")
    team           = optional(string, "devops")
    subnet_newbits = optional(number, 8)
    subnet_azs     = list(string)
    subnets        = optional(map(list(number)), {})
    nat = optional(object({
      create = optional(bool, true)
      per_az = optional(bool, false)
      subnet = optional(string, "")
    }), {})
  })

  validation { # env 값이 develop, staging, rc, production 중 하나인가?
    condition     = contains(["develop", "staging", "rc", "production"], var.attribute.env)
    error_message = "env 값이 develop, staging, rc, production 중 하나여야 합니다."
  }

  validation {
    condition     = alltrue([for az in var.attribute.subnet_azs : can(regex("^[a-zA-Z]$", az))])
    error_message = "subnet_azs 값이 알파벳으로 구성되어야 합니다."
  }

  validation {
    condition     = length(flatten([for k, v in var.attibute.subnets : v])) == length(distinct(flatten([for k, v in var.attibute.subnets : v])))
    error_message = "한 VPC 내에서 서브네팅을 위한 netnum 은 겹칠 수 없습니다."
  }

  validation {
    condition     = alltrue([for k, v in var.attribute.subnets : length(v) <= length(var.attribute.subnet_azs)])
    error_message = "각 subnet 의 netnum list 의 길이는 subnet_azs 의 길이 (${length(var.attribute.subnet_azs)})보다 작거나 같아야 합니다."
  }

  validation {
    condition     = !(var.attribute.nat.create && !contains([for k, v in var.attribute.subnets : k], var.attribute.nat.subnet))
    error_message = "[${local.vpc_name} VPC] nat.subnet 이름은 subnets 에 기재된 항목 중 하나여야 한다."
  }
}

variable "tags" {
  description = "모든 리소스에 적용될 태그 (map)"
  type        = map(string)
  default     = {}
}
