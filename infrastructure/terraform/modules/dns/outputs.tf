output "zone_id" {
  description = "Route53 zone ID"
  value       = local.zone_id
}

output "zone_name" {
  description = "Route53 zone name"
  value       = local.root_domain
}

output "nameservers" {
  description = "Route53 zone nameservers"
  value       = var.create_zone ? aws_route53_zone.primary[0].name_servers : data.aws_route53_zone.existing[0].name_servers
}

output "certificate_arn" {
  description = "ARN of the ACM certificate"
  value       = aws_acm_certificate.primary.arn
}

output "certificate_domain" {
  description = "Domain name of the ACM certificate"
  value       = aws_acm_certificate.primary.domain_name
}

output "certificate_validation_domains" {
  description = "List of domain validation options"
  value       = aws_acm_certificate.primary.domain_validation_options
}

output "api_domain" {
  description = "API domain name"
  value       = local.api_domain
}

output "www_domain" {
  description = "WWW domain name"
  value       = local.www_domain
}

output "streaming_domain" {
  description = "Streaming domain name"
  value       = local.streaming_domain
}

output "health_check_api_id" {
  description = "ID of the API health check"
  value       = aws_route53_health_check.api.id
}

output "health_check_streaming_id" {
  description = "ID of the streaming health check"
  value       = aws_route53_health_check.streaming.id
}
