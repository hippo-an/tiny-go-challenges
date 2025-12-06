variable "name" {
  description = "VPC peering name"
  type        = string
}
variable "requester_vpc_id" {
  description = "requester vpc id"
  type        = string
}
variable "accepter_vpc_id" {
  description = "accepter vpc id"
  type        = string
}
variable "tags" {
  description = "tag for all resources"
  type        = map(string)
  default     = {}
}
