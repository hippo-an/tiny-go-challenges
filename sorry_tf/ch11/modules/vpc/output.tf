output "vpc_name" {
  value       = local.vpc_name
  description = "VPC 이름"
}

output "vpc_id" {
  value       = local.vpc_id
  description = "VPC ID"
}

output "vpc_cidr" {
  value       = local.vpc_cidr
  description = "VPC CIDR"
}

output "subnet_ids" {
  value       = local.subnet_ids
  description = "서브넷 ID"
}

output "subnet_ids_with_az" {
  value       = local.subnet_ids_with_az
  description = "AZ 별 서브넷 ID"
}

output "public_rt_id" {
  value       = local.public_rt
  description = "퍼블릭 라우팅 테이블 ID"
}

output "private_rt_id" {
  value       = local.private_rts
  description = "프라이빗 라우팅 테이블 ID"
}

output "igw_id" {
  value       = try(aws_internet_gateway.this[0].id, null)
  description = "인터넷 게이트웨이 ID"
}

output "nat_ids" {
  value = {
    for k in local.nat_set : k => aws_nat_gateway.this[k].id
  }
  description = "NAT 게이트웨이 ID"
}


