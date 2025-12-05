output "sg_id" {
  description = "보안 그룹 id (map)"
  value = {
    for k, v in aws_security_group.this : k => v.id
  }

}
