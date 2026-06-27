locals {
  name_prefix    = "${var.project_name}-${var.environment}"
  domain_aliases = compact([
    var.domain_name,
    "www.${var.domain_name}",
    "cdn.${var.domain_name}",
    "media.${var.domain_name}",
    var.environment == "production" ? "api.${var.domain_name}" : null,
    var.environment == "production" ? "streaming.${var.domain_name}" : null,
  ])
}

resource "aws_cloudfront_origin_access_control" "media" {
  name                              = "${local.name_prefix}-media-oac"
  description                       = "OAC for media bucket - ${var.environment}"
  origin_access_control_origin_type = "s3"
  signing_behavior                  = "always"
  signing_protocol                  = "sigv4"
}

resource "aws_cloudfront_distribution" "main" {
  enabled             = true
  is_ipv6_enabled     = true
  comment             = "CloudFront distribution for ${var.project_name} - ${var.environment}"
  default_root_object = "index.html"
  price_class         = var.price_class

  aliases = local.domain_aliases

  viewer_certificate {
    acm_certificate_arn      = var.certificate_arn
    ssl_support_method       = "sni-only"
    minimum_protocol_version = "TLSv1.2_2021"
  }

  origin {
    domain_name              = var.alb_domain
    origin_id                = "ALB-Origin"
    origin_shield {
      enabled              = var.environment == "production"
      origin_shield_region = var.environment == "production" ? "us-east-1" : null
    }
    custom_origin_config {
      http_port              = 80
      https_port             = 443
      origin_protocol_policy = "https-only"
      origin_ssl_protocols   = ["TLSv1.2"]
    }
  }

  origin {
    domain_name              = "${var.s3_media_bucket}.s3.amazonaws.com"
    origin_id                = "S3-Media"
    origin_access_control_id = aws_cloudfront_origin_access_control.media.id
  }

  default_cache_behavior {
    target_origin_id       = "ALB-Origin"
    viewer_protocol_policy = "redirect-to-https"
    allowed_methods        = ["DELETE", "GET", "HEAD", "OPTIONS", "PATCH", "POST", "PUT"]
    cached_methods         = ["GET", "HEAD"]
    compress               = true

    forwarded_values {
      query_string = true
      headers      = ["Origin", "Authorization", "CloudFront-Forwarded-Proto", "Host"]
      cookies {
        forward = "all"
      }
    }

    function_association {
      event_type   = "viewer-request"
      function_arn = aws_cloudfront_function.auth_rewrite.arn
    }

    min_ttl     = 0
    default_ttl = 60
    max_ttl     = 3600
  }

  ordered_cache_behavior {
    path_pattern           = "/static/*"
    target_origin_id       = "S3-Media"
    viewer_protocol_policy = "redirect-to-https"
    allowed_methods        = ["GET", "HEAD", "OPTIONS"]
    cached_methods         = ["GET", "HEAD"]
    compress               = true

    forwarded_values {
      query_string = false
      headers      = ["Origin"]
      cookies {
        forward = "none"
      }
    }

    min_ttl     = 0
    default_ttl = 86400
    max_ttl     = 31536000
  }

  ordered_cache_behavior {
    path_pattern           = "/media/*"
    target_origin_id       = "S3-Media"
    viewer_protocol_policy = "redirect-to-https"
    allowed_methods        = ["GET", "HEAD", "OPTIONS"]
    cached_methods         = ["GET", "HEAD"]
    compress               = true

    forwarded_values {
      query_string = false
      headers      = ["Origin"]
      cookies {
        forward = "none"
      }
    }

    min_ttl     = 0
    default_ttl = 3600
    max_ttl     = 86400
  }

  ordered_cache_behavior {
    path_pattern           = "/thumbnails/*"
    target_origin_id       = "S3-Media"
    viewer_protocol_policy = "redirect-to-https"
    allowed_methods        = ["GET", "HEAD", "OPTIONS"]
    cached_methods         = ["GET", "HEAD"]
    compress               = true

    forwarded_values {
      query_string = false
      headers      = ["Origin"]
      cookies {
        forward = "none"
      }
    }

    min_ttl     = 0
    default_ttl = 86400
    max_ttl     = 604800
  }

  ordered_cache_behavior {
    path_pattern           = "/health"
    target_origin_id       = "ALB-Origin"
    viewer_protocol_policy = "https-only"
    allowed_methods        = ["GET", "HEAD"]
    cached_methods         = ["GET", "HEAD"]
    compress               = true

    forwarded_values {
      query_string = false
      headers      = []
      cookies {
        forward = "none"
      }
    }

    min_ttl     = 0
    default_ttl = 0
    max_ttl     = 0
  }

  custom_error_response {
    error_code         = 403
    response_code      = 404
    response_page_path = "/index.html"
  }

  custom_error_response {
    error_code         = 404
    response_code      = 404
    response_page_path = "/index.html"
  }

  custom_error_response {
    error_code         = 500
    response_code      = 500
    response_page_path = "/index.html"
  }

  custom_error_response {
    error_code         = 502
    response_code      = 502
    response_page_path = "/index.html"
  }

  custom_error_response {
    error_code         = 503
    response_code      = 503
    response_page_path = "/index.html"
  }

  geo_restriction {
    restriction_type = var.geo_restriction_type
    locations        = var.geo_restriction_locations
  }

  logging_config {
    include_cookies = false
    bucket          = "${var.s3_media_bucket}.s3.amazonaws.com"
    prefix          = "cloudfront-logs/"
  }

  tags = {
    Name        = "${local.name_prefix}-cloudfront"
    Environment = var.environment
  }
}

resource "aws_cloudfront_function" "auth_rewrite" {
  name    = "${local.name_prefix}-auth-rewrite"
  runtime = "cloudfront-js-2.0"
  comment = "Rewrite auth headers for the Spark platform"
  code    = <<-EOT
function handler(event) {
    var request = event.request;
    var uri = request.uri;

    if (uri.startsWith('/api/') || uri.startsWith('/streaming/')) {
        request.headers['x-forwarded-host'] = { value: request.headers.host.value };
        request.headers['x-original-uri'] = { value: uri };
    }

    if (uri === '/' || uri === '') {
        request.uri = '/index.html';
    }

    return request;
}
EOT
}

resource "aws_wafv2_web_acl" "cloudfront" {
  name        = "${local.name_prefix}-waf"
  description = "WAF ACL for CloudFront distribution - ${var.environment}"
  scope       = "CLOUDFRONT"

  default_action {
    allow {}
  }

  rule {
    name     = "AWSRateBasedRule"
    priority = 1

    action {
      block {}
    }

    statement {
      rate_based_statement {
        limit              = var.environment == "production" ? 10000 : 5000
        aggregate_key_type = "IP"
      }
    }

    visibility_config {
      cloudwatch_metrics_enabled = true
      metric_name               = "${local.name_prefix}-rate-limit"
      sampled_requests_enabled   = true
    }
  }

  rule {
    name     = "AWS-AWSManagedRulesCommonRuleSet"
    priority = 2

    override_action {
      none {}
    }

    statement {
      managed_rule_group_statement {
        name        = "AWSManagedRulesCommonRuleSet"
        vendor_name = "AWS"
      }
    }

    visibility_config {
      cloudwatch_metrics_enabled = true
      metric_name               = "${local.name_prefix}-common-rules"
      sampled_requests_enabled   = true
    }
  }

  rule {
    name     = "AWS-AWSManagedRulesSQLiRuleSet"
    priority = 3

    override_action {
      none {}
    }

    statement {
      managed_rule_group_statement {
        name        = "AWSManagedRulesSQLiRuleSet"
        vendor_name = "AWS"
      }
    }

    visibility_config {
      cloudwatch_metrics_enabled = true
      metric_name               = "${local.name_prefix}-sqli-rules"
      sampled_requests_enabled   = true
    }
  }

  rule {
    name     = "AWS-AWSManagedRulesKnownBadInputsRuleSet"
    priority = 4

    override_action {
      none {}
    }

    statement {
      managed_rule_group_statement {
        name        = "AWSManagedRulesKnownBadInputsRuleSet"
        vendor_name = "AWS"
      }
    }

    visibility_config {
      cloudwatch_metrics_enabled = true
      metric_name               = "${local.name_prefix}-bad-inputs"
      sampled_requests_enabled   = true
    }
  }

  rule {
    name     = "AWS-AWSManagedRulesAmazonIpReputationList"
    priority = 5

    override_action {
      none {}
    }

    statement {
      managed_rule_group_statement {
        name        = "AWSManagedRulesAmazonIpReputationList"
        vendor_name = "AWS"
      }
    }

    visibility_config {
      cloudwatch_metrics_enabled = true
      metric_name               = "${local.name_prefix}-ip-reputation"
      sampled_requests_enabled   = true
    }
  }

  rule {
    name     = "AWS-AWSManagedRulesAnonymousIpList"
    priority = 6

    override_action {
      none {}
    }

    statement {
      managed_rule_group_statement {
        name        = "AWSManagedRulesAnonymousIpList"
        vendor_name = "AWS"
      }
    }

    visibility_config {
      cloudwatch_metrics_enabled = true
      metric_name               = "${local.name_prefix}-anonymous-ip"
      sampled_requests_enabled   = true
    }
  }

  visibility_config {
    cloudwatch_metrics_enabled = true
    metric_name               = "${local.name_prefix}-waf"
    sampled_requests_enabled   = true
  }

  tags = {
    Name        = "${local.name_prefix}-waf"
    Environment = var.environment
  }
}

resource "aws_wafv2_web_acl_association" "cloudfront" {
  resource_arn = aws_cloudfront_distribution.main.arn
  web_acl_arn  = aws_wafv2_web_acl.cloudfront.arn
}

resource "aws_cloudfront_monitoring_subscription" "main" {
  distribution_id = aws_cloudfront_distribution.main.id

  monitoring_subscription {
    realtime_metrics_subscription_config {
      realtime_metrics_subscription_status = "Enabled"
    }
  }
}
