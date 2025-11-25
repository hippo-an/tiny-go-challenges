locals {
    service_name = "web-app"
    common_tags = {
        Service = local.service_name
        ManagedBy = "Terraform"
    }
}

output "output_name" {
    value = <value>
    description = <description>
    sensitive = <true_or_false>
}