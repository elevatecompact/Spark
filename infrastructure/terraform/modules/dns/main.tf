locals {
  name_prefix    = "${var.project_name}-${var.environment}"
  root_domain    = var.domain_name
  api_domain     = "api.${var.domain_name}"
  www_domain     = "www.${var.domain_name}"
  streaming_domain = "streaming.${var.domain_name}"
  cdn_domain     = "cdn.${var.domain_name}"
  media_domain   = "media.${var.domain_name}"
}

resource "aws_route53_zone" "primary" {
  count = var.create_zone ? 1 : 0

  name = local.root_domain

  tags = {
    Name        = "${local.name_prefix}-zone"
    Environment = var.environment
  }
}

data "aws_route53_zone" "existing" {
  count = var.create_zone ? 0 : 1

  name         = local.root_domain
  private_zone = false
}

locals {
  zone_id = var.create_zone ? aws_route53_zone.primary[0].zone_id : data.aws_route53_zone.existing[0].zone_id
}

resource "aws_route53_record" "api" {
  count   = var.alb_domain != "" ? 1 : 0
  zone_id = local.zone_id
  name    = local.api_domain
  type    = "A"

  alias {
    name                   = var.alb_domain
    zone_id                = var.alb_zone_id
    evaluate_target_health = true
  }
}

resource "aws_route53_record" "www" {
  count   = var.cloudfront_domain != "" ? 1 : 0
  zone_id = local.zone_id
  name    = local.www_domain
  type    = "A"

  alias {
    name                   = var.cloudfront_domain
    zone_id                = "Z2FDTNDATAQYW2"
    evaluate_target_health = false
  }
}

resource "aws_route53_record" "streaming" {
  count   = var.alb_domain != "" ? 1 : 0
  zone_id = local.zone_id
  name    = local.streaming_domain
  type    = "A"

  alias {
    name                   = var.alb_domain
    zone_id                = var.alb_zone_id
    evaluate_target_health = true
  }
}

resource "aws_route53_record" "root" {
  count   = var.cloudfront_domain != "" ? 1 : 0
  zone_id = local.zone_id
  name    = local.root_domain
  type    = "A"

  alias {
    name                   = var.cloudfront_domain
    zone_id                = "Z2FDTNDATAQYW2"
    evaluate_target_health = false
  }
}

resource "aws_route53_record" "cdn" {
  count   = var.cloudfront_domain != "" ? 1 : 0
  zone_id = local.zone_id
  name    = local.cdn_domain
  type    = "CNAME"
  ttl     = 300
  records = [var.cloudfront_domain]
}

resource "aws_route53_record" "media" {
  count   = var.cloudfront_domain != "" ? 1 : 0
  zone_id = local.zone_id
  name    = local.media_domain
  type    = "CNAME"
  ttl     = 300
  records = [var.cloudfront_domain]
}

resource "aws_route53_record" "mx" {
  zone_id = local.zone_id
  name    = local.root_domain
  type    = "MX"
  ttl     = 300
  records = [
    "10 inbound-smtp.${var.region}.amazonaws.com",
    "10 feedback-smtp.${var.region}.amazonaws.com"
  ]
}

resource "aws_route53_record" "spf" {
  zone_id = local.zone_id
  name    = local.root_domain
  type    = "TXT"
  ttl     = 300
  records = [
    "v=spf1 include:amazonses.com include:_spf.google.com ~all"
  ]
}

resource "aws_route53_record" "dkim" {
  zone_id = local.zone_id
  name    = "${local.dkim_selector}._domainkey.${local.root_domain}"
  type    = "CNAME"
  ttl     = 300
  records = ["${local.dkim_selector}.dkim.amazonses.com"]
}

resource "aws_route53_record" "dmarc" {
  zone_id = local.zone_id
  name    = "_dmarc.${local.root_domain}"
  type    = "TXT"
  ttl     = 300
  records = [
    "v=DMARC1; p=quarantine; pct=100; rua=mailto:dmarc-reports@${local.root_domain}"
  ]
}

resource "aws_route53_record" "domainkey" {
  zone_id = local.zone_id
  name    = "default._domainkey.${local.root_domain}"
  type    = "TXT"
  ttl     = 300
  records = [
    "v=DKIM1; h=sha256; k=rsa; p=MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDL4g7gKEb4+9vJ+5vq8a2e5vL6f8g7J0v3c5b4d5e6f7g8h9i0j1k2l3m4n5o6p7q8r9s0t1u2v3w4x5y6z7A8B9C0D1E2F3G4H5I6J7K8L9M0N1O2P3Q4R5S6T7U8V9W0X1Y2Z3A4B5C6D7E8F9G0H1I2J3K4L5M6N7O8P9Q0R1S2T3U4V5W6X7Y8Z9A0B1C2D3E4F5G6H7I8J9K0L1M2N3O4P5Q6R7S8T9UA0B1C2D3E4F5G6H7I8J9K0L1M2N3O4P5Q6R7S8T9U"
  ]
}

resource "aws_route53_health_check" "api" {
  fqdn              = local.api_domain
  port              = 443
  type              = "HTTPS"
  resource_path     = "/health"
  failure_threshold = 3
  request_interval  = 30

  tags = {
    Name        = "${local.name_prefix}-api-health-check"
    Environment = var.environment
  }
}

resource "aws_route53_health_check" "streaming" {
  fqdn              = local.streaming_domain
  port              = 443
  type              = "HTTPS"
  resource_path     = "/health"
  failure_threshold = 3
  request_interval  = 30

  tags = {
    Name        = "${local.name_prefix}-streaming-health-check"
    Environment = var.environment
  }
}

resource "aws_acm_certificate" "primary" {
  domain_name       = local.root_domain
  subject_alternative_names = [
    local.api_domain,
    local.www_domain,
    local.streaming_domain,
    local.cdn_domain,
    local.media_domain,
    "*.${local.root_domain}"
  ]
  validation_method = "DNS"

  tags = {
    Name        = "${local.name_prefix}-certificate"
    Environment = var.environment
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_route53_record" "cert_validation" {
  for_each = {
    for dvo in aws_acm_certificate.primary.domain_validation_options : dvo.domain_name => {
      name   = dvo.resource_record_name
      record = dvo.resource_record_value
      type   = dvo.resource_record_type
    }
  }

  allow_overwrite = true
  name            = each.value.name
  records         = [each.value.record]
  ttl             = 60
  type            = each.value.type
  zone_id         = local.zone_id
}

resource "aws_acm_certificate_validation" "primary" {
  certificate_arn         = aws_acm_certificate.primary.arn
  validation_record_fqdns = [for record in aws_route53_record.cert_validation : record.fqdn]
}
