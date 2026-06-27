output "cloudfront_id" {
  description = "CloudFront distribution ID"
  value       = aws_cloudfront_distribution.main.id
}

output "cloudfront_domain" {
  description = "CloudFront distribution domain name"
  value       = aws_cloudfront_distribution.main.domain_name
}

output "cloudfront_arn" {
  description = "ARN of the CloudFront distribution"
  value       = aws_cloudfront_distribution.main.arn
}

output "cloudfront_hosted_zone_id" {
  description = "CloudFront hosted zone ID for Route53 aliases"
  value       = aws_cloudfront_distribution.main.hosted_zone_id
}

output "waf_acl_arn" {
  description = "ARN of the WAF ACL"
  value       = aws_wafv2_web_acl.cloudfront.arn
}

output "waf_acl_id" {
  description = "ID of the WAF ACL"
  value       = aws_wafv2_web_acl.cloudfront.id
}

output "origin_access_control_id" {
  description = "ID of the S3 origin access control"
  value       = aws_cloudfront_origin_access_control.media.id
}
