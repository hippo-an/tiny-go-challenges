variable "vpc_name" {
  description = "EC2 가 존재할 vpc 이름"
  type        = string
}

variable "vpc_id" {
  description = "EC2 가 존재할 vpc id"
  type        = string
}

variable "subnet_id_map" {
  description = "subnet id 맵 데이터"
  type        = map(map(string))
}

variable "sg_id_map" {
  description = "sg id 맵 데이터"
  type        = map(string)
}


variable "tags" {
  description = "모든 리소스에 적용될 태그"
  type        = map(string)
  default     = {}
}

variable "ec2_set" {
  description = "EC2 인스턴스 별 명세 Set"
  type = map(object({
    # required
    env             = string
    team            = string
    service         = string
    ami_id          = string
    instance_type   = string
    subnet          = string
    az              = string
    security_groups = list(string)


    # optional
    os_type    = optional(string, "linux")
    ec2_key    = optional(string)
    ec2_role   = optional(string)
    private_ip = optional(string)

    # root volume
    root_volume = object({
      size = number
      type = optional(string, "gp3")
    })

    # additional volumes
    additional_volume = optional(list(object({
      device = string
      size   = number
      type   = optional(string, "gp3")
      iops   = optional(number)
    })), [])
  }))
}

