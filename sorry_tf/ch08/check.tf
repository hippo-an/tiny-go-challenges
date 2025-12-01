resource "aws_acm_certificate" "validation" {
  domain_name               = "terraform.com"
  validation_method         = "DNS"
  subject_alternative_names = ["terraform.com"]
}

locals {
  acm_validate_value = tolist(aws_acm_certificate.validation.domain_validation_options)[0]
}

resource "time_sleep" "wait" {
  depends_on      = [aws_acm_certificate.validation]
  create_duration = "30s"
}

check "validation" {
  data "aws_acm_certificate" "this" {
    domain     = aws_acm_certificate.validation.domain_name
    statuses   = ["VALIDATION_TIME_OUT", "PENDING_VALIDATION", "EXPIRED", "INACTIVE", "REVOKED", "ISSUED", "FAILED"]
    depends_on = [time_sleep.wait]
  }

  assert {
    condition     = data.aws_acm_certificate.this.status != "PENDING_VALIDATION"
    error_message = "ACM certificate validation failed"
  }
}

check "validation_record" {
  data "dns_cname_record_set" "validation" {
    host       = local.acm_validate_value.resource_record_name
    depends_on = [time_sleep.wait]
  }

  assert {
    condition     = data.dns_cname_record_set.validation.cname == local.acm_validate_value.resource_record_name
    error_message = "ACM certificate validation record is not a CNAME"
  }
}
