variable "vpc_name" {
  description = "보안 그룹이 존재할 VPC 의 이름"
  type        = string
}
variable "vpc_id" {
  type        = string
  description = "보안 그룹이 존재할 VPC 의 ID"
}

variable "sg_set" {
  description = "보안 그룹별 인바운드 ruleset"
  type = map(object({
    desc      = optional(string, "")
    protocol  = string
    from_port = number
    to_port   = number
    source    = string

  }))
}
variable "tags" {
  description = "모든 리소스에 적용될 태그 (map)"
  type        = map(string)
  default     = {}
}
