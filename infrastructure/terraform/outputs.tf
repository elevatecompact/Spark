output "vpc_id" {
  description = "ID of the VPC"
  value       = module.networking.vpc_id
}

output "vpc_cidr" {
  description = "CIDR block of the VPC"
  value       = module.networking.vpc_cidr
}

output "cluster_endpoint" {
  description = "EKS cluster API endpoint"
  value       = module.kubernetes.cluster_endpoint
}

output "cluster_name" {
  description = "EKS cluster name"
  value       = module.kubernetes.cluster_name
}

output "database_endpoint" {
  description = "RDS database endpoint"
  value       = module.database.endpoint
  sensitive   = true
}

output "database_name" {
  description = "RDS database name"
  value       = module.database.database_name
}

output "dns_nameservers" {
  description = "Route53 zone nameservers"
  value       = module.dns.nameservers
}

output "cloudfront_domain" {
  description = "CloudFront distribution domain name"
  value       = module.cdn.cloudfront_domain
}

output "media_bucket" {
  description = "S3 media bucket name"
  value       = module.storage.media_bucket_name
}

output "environment" {
  description = "Current deployment environment"
  value       = var.environment
}
